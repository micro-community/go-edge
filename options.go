package edge

import (
	"context"
	"crypto/tls"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

// Options  of edge node serivices
type Options struct {
	Name      string
	Version   string
	ID        string
	Metadata  map[string]string
	Address   string
	Advertise string
	Action    func(*cli.Context)
	Flags     []cli.Flag

	Extractor      DataExtractor
	IsTCPTransport bool //indentity tcpsession with node or not

	Context context.Context // Alternative Options
	Service micro.Service

	Secure      bool
	TLSConfig   *tls.Config
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:           DefaultName,
		Version:        DefaultVersion,
		ID:             DefaultID,
		Address:        DefaultAddress,
		Service:        micro.NewService(),
		Context:        context.TODO(),
		Extractor:      DefaultExtractor,
		IsTCPTransport: true,
	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}

//Name of Server name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Icon specifies an icon url to load in the UI
func Icon(ico string) Option {
	return func(o *Options) {
		if o.Metadata == nil {
			o.Metadata = make(map[string]string)
		}
		o.Metadata["icon"] = ico
	}
}

//ID Unique server id
func ID(id string) Option {
	return func(o *Options) {
		o.ID = id
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Metadata associated with the service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Metadata = md
	}
}

// Address to bind to - host:port
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

//Advertise for discovery - host:port
func Advertise(a string) Option {
	return func(o *Options) {
		o.Advertise = a
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// //Server of edge service
// func Server(srv *edge.Server) Option {
// 	return func(o *Options) {
// 		o.Server = srv
// 	}
// }

// MicroService sets the micro.Service used internally
func MicroService(s micro.Service) Option {
	return func(o *Options) {
		o.Service = s
	}
}

// Flags sets the command flags.
func Flags(flags ...cli.Flag) Option {
	return func(o *Options) {
		o.Flags = append(o.Flags, flags...)
	}
}

// Action sets the command action.
func Action(a func(*cli.Context)) Option {
	return func(o *Options) {
		o.Action = a
	}
}

// BeforeStart is executed before the server starts.
func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

// BeforeStop is executed before the server stops.
func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

// AfterStart is executed after server start.
func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

// AfterStop is executed after server stop.
func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}

// Secure Use secure communication. If TLSConfig is not specified we use InsecureSkipVerify and generate a self signed cert
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// TLSConfig to be used for the transport.
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

//Extractor edge message
func Extractor(e DataExtractor) Option {
	return func(o *Options) {
		o.Extractor = e
	}
}

//WithTransport of
func WithTransport(isTCP bool) Option {
	return func(o *Options) {
		o.IsTCPTransport = isTCP
	}
}
