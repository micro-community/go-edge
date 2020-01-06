package experiment

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	edge "github.com/micro-community/x-edge"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	var (
		beforeStartCalled bool
		afterStartCalled  bool
		beforeStopCalled  bool
		afterStopCalled   bool
		str               = `<html><body><h1>Hello World</h1></body></html>`
		fn                = func(w edge.ResponseWriter, r *edge.Request) { fmt.Fprint(w, str) }
	)

	beforeStart := func() error {
		beforeStartCalled = true
		return nil
	}

	afterStart := func() error {
		afterStartCalled = true
		return nil
	}

	beforeStop := func() error {
		beforeStopCalled = true
		return nil
	}

	afterStop := func() error {
		afterStopCalled = true
		return nil
	}

	service := NewService(
		Name("go.micro.web.test"),
		BeforeStart(beforeStart),
		AfterStart(afterStart),
		BeforeStop(beforeStop),
		AfterStop(afterStop),
	)

	service.HandleFunc("/", fn)

	errCh := make(chan error, 1)
	go func() {
		errCh <- service.Run()
		close(errCh)
	}()

	eventually(func() bool {
		var err error
		//		s, err = reg.GetService("go.micro.web.test")
		return err == nil
	}, t.Fatal)

	if have, want := len(s), 1; have != want {
		t.Fatalf("Expected %d but got %d services", want, have)
	}

	rsp, err := edge.Get(fmt.Sprintf("edge://%s", s[0].Nodes[0].Address))
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != str {
		t.Errorf("Expected %s got %s", str, string(b))
	}

	callbackTests := []struct {
		subject string
		have    interface{}
	}{
		{"beforeStartCalled", beforeStartCalled},
		{"afterStartCalled", afterStartCalled},
	}

	for _, tt := range callbackTests {
		if tt.have != true {
			t.Errorf("unexpected %s: want true, have false", tt.subject)
		}
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("service.Run():%v", err)
		}
	case <-time.After(time.Duration(time.Second)):
		t.Logf("service.Run() survived a client request without an error")
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)

	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	<-ch

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("service.Run():%v", err)
		} else {
			t.Log("service.Run() nil return on syscall.SIGTERM")
		}
	case <-time.After(time.Duration(time.Second)):
		t.Logf("service.Run() survived a client request without an error")
	}

	callbackTests = []struct {
		subject string
		have    interface{}
	}{
		{"beforeStopCalled", beforeStopCalled},
		{"afterStopCalled", afterStopCalled},
	}

	for _, tt := range callbackTests {
		if tt.have != true {
			t.Errorf("unexpected %s: want true, have false", tt.subject)
		}
	}

}

func TestOptions(t *testing.T) {
	var (
		name      = "service-name"
		id        = "service-id"
		version   = "service-version"
		address   = "service-addr"
		advertise = "service-adv"
		handler   = edge.NewServeMux()
		metadata  = map[string]string{"key": "val"}
		secure    = true
	)

	service := NewService(
		Name(name),
		ID(id),
		Version(version),
		Address(address),
		Advertise(advertise),
		Handler(handler),
		Metadata(metadata),
		Secure(secure),
	)

	opts := service.Options()

	tests := []struct {
		subject string
		want    interface{}
		have    interface{}
	}{
		{"name", name, opts.Name},
		{"version", version, opts.Version},
		{"id", id, opts.ID},
		{"address", address, opts.Address},
		{"advertise", advertise, opts.Advertise},
		{"handler", handler, opts.Handler},
		{"metadata", metadata["key"], opts.Metadata["key"]},
		{"secure", secure, opts.Secure},
	}

	for _, tc := range tests {
		if tc.want != tc.have {
			t.Errorf("unexpected %s: want %v, have %v", tc.subject, tc.want, tc.have)
		}
	}
}

func eventually(pass func() bool, fail func(...interface{})) {
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()

	timeout := time.After(time.Second)

	for {
		select {
		case <-timeout:
			fail("timed out")
			return
		case <-tick.C:
			if pass() {
				return
			}
		}
	}
}

func TestTLS(t *testing.T) {
	var (
		str    = `<html><body><h1>Hello World</h1></body></html>`
		fn     = func(w edge.ResponseWriter, r *edge.Request) { fmt.Fprint(w, str) }
		secure = true
	)

	service := NewService(
		Name("go.micro.web.test"),
		Secure(secure),
	)

	service.HandleFunc("/", fn)

	errCh := make(chan error, 1)
	go func() {
		errCh <- service.Run()
		close(errCh)
	}()

	eventually(func() bool {
		var err error
		//		s, err = reg.GetService("go.micro.web.test")
		return err == nil
	}, t.Fatal)

	if have, want := len(s), 1; have != want {
		t.Fatalf("Expected %d but got %d services", want, have)
	}

	tr := &edge.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &edge.Client{Transport: tr}
	rsp, err := client.Get(fmt.Sprintf("edges://%s", s[0].Nodes[0].Address))
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != str {
		t.Errorf("Expected %s got %s", str, string(b))
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("service.Run():%v", err)
		}
	case <-time.After(time.Duration(time.Second)):
		t.Logf("service.Run() survived a client request without an error")
	}

}
