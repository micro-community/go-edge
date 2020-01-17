package edge

import (
	"errors"

	nclient "github.com/micro-community/x-edge/node/client"
	nserver "github.com/micro-community/x-edge/node/server"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/service"
)

//Default Config
var (
	// For serving node connection
	DefaultName    = "x-edge-node-srv"
	DefaultAddress = ":8000"

	DefaultExtractor = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return -1, nil, errors.New("No Extractor Defined")
	}
)

type edgeService struct {
	opts service.Options
}

func newService(opts ...service.Option) service.Service {
	options := service.NewOptions(opts...)
	return &edgeService{
		opts: options,
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

func (s *edgeService) Run() error {
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
