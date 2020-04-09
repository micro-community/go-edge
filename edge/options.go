package edge

import (
	"context"
	"crypto/tls"
	"time"

	nts "github.com/micro-community/x-edge/node/transport"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/transport"
)

//Options for edge Service
type Options struct {
	//	service.Options //inherit from service
	Name      string
	Version   string
	ID        string
	Metadata  map[string]string
	Address   string
	Advertise string

	Auth   auth.Auth
	Client client.Client
	Server server.Server

	Transport transport.Transport
	Action    func(*cli.Context)
	Flags     []cli.Flag

	RegisterTTL      time.Duration
	RegisterInterval time.Duration

	//	Handler http.Handler

	// Alternative Options
	Context context.Context

	Secure      bool
	TLSConfig   *tls.Config
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	Signal    bool
	Namespace string
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:      DefaultName,
		Version:   DefaultVersion,
		ID:        DefaultID,
		Address:   DefaultAddress,
		Server:    DefaultServer,
		Transport: DefaultTransport,
		Context:   context.TODO(),
		Signal:    true,
	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Name Server name
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

// ID Unique server id
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

// Advertise The address to advertise for discovery - host:port
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

// RegisterTTL the service with a TTL
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = t
	}
}

// RegisterInterval Register the service with at interval
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = t
	}
}

// // Handler binding
// func Handler(h http.Handler) Option {
// 	return func(o *Options) {
// 		o.Handler = h
// 	}
// }

// Server to use a customer Server
func Server(srv server.Server) Option {
	return func(o *Options) {
		o.Server = srv
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

// HandleSignal toggles automatic installation of the signal handler that
// traps TERM, INT, and QUIT.  Users of this feature to disable the signal
// handler, should control liveness of the service through the context.
func HandleSignal(b bool) Option {
	return func(o *Options) {
		o.Signal = b
	}
}

// Options  of edge node serivices

//WithExtractor edge message
func WithExtractor(de nts.DataExtractor) Option {
	return func(o *Options) {
		o.Transport.Init(nts.WithExtractor(de))
	}
}

// Transport sets the transport for the server
// and the underlying components
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
		// Update Client and Server
		o.Client.Init(client.Transport(t))
		o.Server.Init(server.Transport(t))
	}
}

//Namespace of edge server
func Namespace(n string) Option {
	return func(o *Options) {
		o.Namespace = n
	}
}
