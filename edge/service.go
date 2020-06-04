package edge

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	nts "github.com/micro-community/x-edge/node/transport"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/debug/trace"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/util/wrapper"
)

type service struct {
	opts Options
	sync.Mutex
	running bool
	foredge bool
	exit    chan chan error
}

func newService(opts ...Option) Service {

	options := newOptions(opts...)

	// service name
	serviceName := options.Server.Options().Name

	// TODO: better accessors
	//authFn := func() auth.Auth { return options.Auth }

	// wrap client to inject From-Service header on any calls
	options.Client = wrapper.FromService(serviceName, options.Client)
	options.Client = wrapper.TraceCall(serviceName, trace.DefaultTracer, options.Client)

	return &service{
		opts:    options,
		foredge: true,
	}
}

func (s *service) Name() string {
	return s.opts.Server.Options().Name
}

// Init initialize options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *service) Init(opts ...Option) error {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

	//var authOpts []auth.Option
	var serverOpts []server.Option
	var clientOpts []client.Option

	s.opts.Action = func(ctx *cli.Context) {

		if len(ctx.String("edge_web_address")) > 0 {
			s.opts.Address = ctx.String("edge_address")
			//clientOpts = append(clientOpts, client.WithAddress(ctx.String("edge_address")))
		}
		if len(ctx.String("edge_host")) > 0 {
			s.opts.Host = ctx.String("edge_host")
			serverOpts = append(serverOpts, server.Address(s.opts.Host))

		}
		if len(ctx.String("edge_address")) > 0 {
			serverOpts = append(serverOpts, server.Address(ctx.String("edge_address")))
		}
		if name := ctx.String("edge_transport"); len(name) > 0 && s.opts.Transport.String() != name {
			//to see if we have a setting up from flag or env
			if t, ok := s.opts.Transports[name]; ok {
				s.opts.Transport = t()
				// to remember we have a extractor to set
				s.opts.Transport.Init(nts.WithExtractor(s.opts.Extractor))
				serverOpts = append(serverOpts, server.Transport(s.opts.Transport))
				clientOpts = append(clientOpts, client.Transport(s.opts.Transport))
			}
		}

		//set Opts
		if len(serverOpts) > 0 {
			if err := s.Server().Init(serverOpts...); err != nil {
				logger.Fatalf("Error configuring server: %v", err)
			}
		}
		//set Opts for client
		if len(clientOpts) > 0 {
			if err := s.Client().Init(clientOpts...); err != nil {
				logger.Fatalf("Error configuring client: %v", err)
			}
		}

	}

	return nil
}

func (s *service) String() string {
	return "edge"
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Client() client.Client {
	return s.opts.Client
}

func (s *service) Server() server.Server {
	return s.opts.Server
}

func (s *service) Start() error {

	s.Lock()
	defer s.Unlock()

	if s.running {
		return nil
	}

	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := s.opts.Server.Start(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

//Run edge srv node
func (s *service) Run() error {

	if err := s.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	}

	select {
	// wait on kill signal
	case sig := <-ch:
		if logger.V(logger.InfoLevel, log) {
			log.Infof("Received signal %s", sig)
		}
	// wait on context cancel
	case <-s.opts.Context.Done():
		if logger.V(logger.InfoLevel, log) {
			log.Info("Received context shutdown")
		}
	}

	return s.Stop()
}

func (s *service) Stop() error {

	s.Lock()
	defer s.Unlock()

	if !s.running {
		return nil
	}

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			return err
		}
	}

	ch := make(chan error, 1)
	s.exit <- ch
	s.running = false

	if logger.V(logger.InfoLevel, log) {
		log.Info("Stopping")
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			if chErr := <-ch; chErr != nil {
				return chErr
			}
			return err
		}
	}

	return <-ch
}
