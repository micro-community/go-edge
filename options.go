package edge

import (
	"context"

	"github.com/micro/go-micro/v2"
)

//Options of edge srv
type Options struct {
	Service micro.Service
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
	return opt
}

// MicroService sets the micro.Service used internally
func MicroService(s micro.Service) Option {
	return func(o *Options) {
		o.Service = s
	}
}
