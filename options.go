package edge

import (
	"context"
	"time"

	nedge "github.com/micro-community/x-edge/edge"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
)

//Options of edge app
type Options struct {
	//Global Metadata
	Name             string
	Version          string
	ID               string
	Metadata         map[string]string
	Address          string
	Advertise        string
	Action           func(*cli.Context)
	Flags            []cli.Flag
	RegisterTTL      time.Duration
	RegisterInterval time.Duration
	// RegisterCheck runs a check function before registering the service
	RegisterCheck func(context.Context) error
	Registry      registry.Registry
	Service       micro.Service

	Edge nedge.Service

	//for edge server
	EdgeTransport string
	EdgeHost      string

	// Alternative Options
	Context context.Context
}

//Option of edge app
type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		Service: micro.NewService(),
		Context: context.TODO(),
	}

	for _, o := range opts {
		o(&opt)
	}

	if opt.Service == nil {
		opt.Service = micro.NewService()
	}
	if opt.Edge == nil {
		opt.Edge = nedge.NewServer()
	}
	return opt
}

// MicroService sets the micro.Service for internal communication
func MicroService(s micro.Service) Option {
	return func(o *Options) {
		o.Service = s
	}
}

// MicroEdge sets the edge.Service for end/controller/gw communication
func MicroEdge(e nedge.Service) Option {
	return func(o *Options) {
		o.Edge = e
	}
}
