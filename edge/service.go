package edge

import (
	nclient "github.com/micro-community/x-edge/node/client"
	nserver "github.com/micro-community/x-edge/node/server"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
)

type service struct {
	opts    service.Options
	service micro.Service
}

func newService(opts ...Option) Service {
	options := service.NewOptions(opts...)
	return &service{
		opts:    options,
		service: micro.NewService(),
	}
}

func (s *edgeService) Name() string {
	return s.opts.Server.Options().Name
}

// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *edgeService) Init(opts ...service.Option) {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

}

func (s *edgeService) String() string {
	return "edge"
}

func (s *edgeService) Options() service.Options {
	return s.opts
}

func (s *edgeService) Client() client.Client {
	return s.opts.Client
}

func (s *edgeService) Server() server.Server {
	return s.opts.Server
}

func (s *edgeService) Start() error {
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
func (s *edgeService) Run() error {

	log.Init(log.WithFields(map[string]interface{}{"service": "edge"}))

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if err := s.Start(); err != nil {
		return err
	}

	// wait on context cancel
	<-s.opts.Context.Done()

	return s.Stop()
}

func (s *edgeService) Stop() error {
	var gerr error

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	if err := s.opts.Server.Stop(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

// return micro.service
func (s *edgeService) MicroService() micro.Service {
	return s.service
}

// NewService returns a new web.Service
func NewService(opts ...service.Option) service.Service {

	// our grpc client
	c := nclient.NewClient()
	// our grpc server
	s := nserver.NewServer()

	// create options with priority for our opts
	options := []service.Option{
		service.Client(c),
		service.Server(s),
	}

	// append passed in opts
	options = append(opts, opts...)

	// generate and return a service
	return newService(options...)
}
