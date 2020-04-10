package edge

import (
	"context"
	"time"

	nedge "github.com/micro-community/x-edge/edge"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/transport"
)

//Options of edge app
type Options struct {
	//Global Metadata
	Namespace        string
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

	//two service end
	MicroService micro.Service
	Edge         nedge.Service

	//for edge server
	EdgeTransport transport.Transport
	EdgeHost      string

	// Alternative Options
	Context context.Context
}

//Option of edge app
type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		MicroService: micro.NewService(),
		Context:      context.TODO(),
	}

	for _, o := range opts {
		o(&opt)
	}

	if opt.MicroService == nil {
		opt.MicroService = micro.NewService()
	}
	if opt.Edge == nil {
		opt.Edge = nedge.NewServer()
	}
	return opt
}

// MicroService sets the micro.Service for internal communication
func MicroService(s micro.Service) Option {
	return func(o *Options) {
		o.MicroService = s
	}
}

// MicroEdge sets the edge.Service for end/controller/gw communication
func MicroEdge(e nedge.Service) Option {
	return func(o *Options) {
		o.Edge = e
	}
}

//Namespace of edge server
func Namespace(n string) Option {
	return func(o *Options) {
		o.Namespace = n
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// EgTransport of the edge
func EgTransport(et transport.Transport) Option {
	return func(o *Options) {
		o.EdgeTransport = et
	}
}
