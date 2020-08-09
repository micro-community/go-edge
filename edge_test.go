package edge

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/service/mucp"
)

func TestInterfaceEqual(t *testing.T) {

	var ms = micro.NewService() //this is a micro.Service
	var ss = mucp.NewService()  //this is a service.Service

	if reflect.TypeOf(ms).Kind() != reflect.TypeOf(ss).Kind() {
		t.Errorf("different type %v : %v", ms, ss)
	}
	//so they are equal type

}

func TestFunction(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	//r := memory.NewRegistry(memory.Services(test.Data))

	ch := make(chan error, 2)

	go func() {
		fmt.Println("doing sth")
		ch <- nil
		wg.Done()
	}()

	// wait for start
	wg.Wait()

	if err := <-ch; err != nil {
		t.Fatal(err)
	}
}
