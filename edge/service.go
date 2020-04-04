package edge

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/debug/trace"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/util/wrapper"

	_ "github.com/micro-community/x-edge/node/client"
	nserver "github.com/micro-community/x-edge/node/server"
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
	// options.Server is empty, so it is crash
	// we should new options.Server  first
	options.Server = nserver.NewServer()

	serviceName := options.Server.Options().Name

	// TODO: better accessors
	authFn := func() auth.Auth { return options.Auth }

	// wrap client to inject From-Service header on any calls
	options.Client = wrapper.FromService(serviceName, options.Client, authFn)
	options.Client = wrapper.TraceCall(serviceName, trace.DefaultTracer, options.Client)

	return &service{
		opts:    options,
		foredge: true,
	}
}

func (s *service) Name() string {
	return s.opts.Server.Options().Name
}

// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *service) Init(opts ...Option) error {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

	serviceOpts := []micro.Option{}

	if len(s.opts.Flags) > 0 {
		serviceOpts = append(serviceOpts, micro.Flags(s.opts.Flags...))
	}
	serviceOpts = append(serviceOpts, micro.Action(func(ctx *cli.Context) error {

		if name := ctx.String("server_name"); len(name) > 0 {
			s.opts.Name = name
		}

		if ver := ctx.String("server_version"); len(ver) > 0 {
			s.opts.Version = ver
		}

		if id := ctx.String("server_id"); len(id) > 0 {
			s.opts.ID = id
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

		s.opts.Service.Init(serviceOpts...)

		return nil
	}))

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

//Run edge service
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

// return micro.service
func (s *service) MicroService() micro.Service {
	return s.Options().Service
}
