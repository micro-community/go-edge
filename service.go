package edge

import (
	"sync"

	"github.com/micro-community/x-edge/config"
	ecli "github.com/micro-community/x-edge/node/client"
	esrv "github.com/micro-community/x-edge/node/server"
	"github.com/micro-community/x-edge/node/transport/tcp"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/micro/cli"
	"github.com/micro/go-micro/util/log"
)

type service struct {
	opts Options
	//	mux *http.ServeMux
	sync.Mutex
	running bool
	static  bool
	exit    chan chan error
}

func newService(opts ...Option) micro.Service {
	options := newOptions(opts...)
	nodeService := micro.NewService(
		micro.Server(esrv.NewServer()),
		micro.Version(config.BuildVersion()),
		micro.Transport(tcp.NewTransport()),
	)

	options.Service = nodeService

	s := &service{
		opts:   options,
		static: true,
	}

	return s
}

func (s *service) start() error {
	s.Lock()
	defer s.Unlock()

	if s.running {
		return nil
	}

	l, err := s.listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	s.opts.Address = l.Addr().String()

	s.exit = make(chan chan error, 1)
	s.running = true

	go func() {
		ch := <-s.exit
		ch <- l.Close()
	}()

	log.Logf("Listening on %v", l.Addr().String())
	return nil
}

func (s *service) stop() error {
	s.Lock()
	defer s.Unlock()

	if !s.running {
		return nil
	}

	ch := make(chan error, 1)
	s.exit <- ch
	s.running = false
	log.Log("Stopping")

	return <-ch
}

func (s *service) Client() client.Client {
	// rt := mhttp.NewRoundTripper(
	// 	mhttp.WithRegistry(registry.DefaultRegistry),
	// )
	return ecli.NewClient()
}

func (s *service) Handle(pattern string, handler node.Handler) {

	// register the handler
	//	s.mux.Handle(pattern, handler)
}

func (s *service) HandleFunc(pattern string, handler func(node.ResponseWriter, *node.Request)) {
	var seen bool
	// for _, ep := range s.srv.Endpoints {
	// 	if ep.Name == pattern {
	// 		seen = true
	// 		break
	// 	}
	// }
	// if !seen {
	// 	s.srv.Endpoints = append(s.srv.Endpoints, &registry.Endpoint{
	// 		Name: pattern,
	// 	})
	// }

	//s.mux.HandleFunc(pattern, handler)
}

func (s *service) Init(opts ...Option) {
	for _, o := range opts {
		o(&s.opts)
	}

	serviceOpts := []micro.Option{}

	if len(s.opts.Flags) > 0 {
		serviceOpts = append(serviceOpts, micro.Flags(s.opts.Flags...))
	}

	serviceOpts = append(serviceOpts, micro.Action(func(ctx *cli.Context) {

		if name := ctx.String("server_name"); len(name) > 0 {
			s.opts.Name = name
		}

		if ver := ctx.String("server_version"); len(ver) > 0 {
			s.opts.Version = ver
		}

		if addr := ctx.String("server_address"); len(addr) > 0 {
			s.opts.Address = addr
		}

		if adv := ctx.String("server_advertise"); len(adv) > 0 {
			s.opts.Advertise = adv
		}

		if s.opts.Action != nil {
			s.opts.Action(ctx)
		}
	}))

	s.opts.Service.Init(serviceOpts...)
	//	srv := s.genSrv()
	//	srv.Endpoints = s.srv.Endpoints
	//	s.srv = srv
}

func (s *service) Run() error {
	if err := s.start(); err != nil {
		return err
	}

	// start reg loop
	ex := make(chan bool)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	select {
	// wait on kill signal
	case sig := <-ch:
		log.Logf("Received signal %s", sig)
	// wait on context cancel
	case <-s.opts.Context.Done():
		log.Logf("Received context shutdown")
	}

	// exit reg loop
	close(ex)

	return s.stop()
}

// Options returns the options for the given service
func (s *service) Options() Options {
	return s.opts
}

func (s *service) listen(network, addr string) (net.Listener, error) {
	var l net.Listener

	return l, nil
}
