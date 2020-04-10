package edge

import (
	nedge "github.com/micro-community/x-edge/edge"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	log "github.com/micro/go-micro/v2/logger"
)

//basic metadata
var (
	//	Transport     = "udp"
	Host = ":8080"
	//  AppNamespace = "x.edge"
	HeaderPrefix = "x-edge-"
)

//Service of edge srv
type Service interface {
	Name() string
	Init(opts ...Option) error
	Run() error
	String() string
	MicroService() micro.Service
	EdgeService() nedge.Service
}

//edgeApp for edge process
type edgeApp struct {
	opts Options
	// auth service
	auth auth.Auth
}

//NewService return a edge service application
func NewService(opts ...Option) Service {

	options := newOptions(opts...)

	s := &edgeApp{
		opts: options,
	}

	return s

}

func (e *edgeApp) buildGoMicroOption() []micro.Option {

	serviceOpts := []micro.Option{}

	if len(e.opts.Flags) > 0 {
		serviceOpts = append(serviceOpts, micro.Flags(e.opts.Flags...))
	}
	serviceOpts = append(serviceOpts, micro.Action(func(ctx *cli.Context) error {

		if name := ctx.String("server_name"); len(name) > 0 {
			e.opts.Name = name
		}

		if ver := ctx.String("server_version"); len(ver) > 0 {
			e.opts.Version = ver
		}

		if id := ctx.String("server_id"); len(id) > 0 {
			e.opts.ID = id
		}

		if addr := ctx.String("server_address"); len(addr) > 0 {
			e.opts.Address = addr
		}

		if adv := ctx.String("server_advertise"); len(adv) > 0 {
			e.opts.Advertise = adv
		}

		return nil
	}))

	return serviceOpts

}

func (e *edgeApp) Init(opts ...Option) error {
	for _, o := range opts {
		o(&e.opts)
	}

	serviceOpts := e.buildGoMicroOption()

	edgeOptions := []nedge.Option{}

	serviceOpts = append(serviceOpts, micro.Action(func(ctx *cli.Context) error {
		if len(ctx.String("edge_web_address")) > 0 {
			edgeOptions = append(edgeOptions, nedge.Address(ctx.String("edge_address")))
		}
		if len(ctx.String("edge_host")) > 0 {
			edgeOptions = append(edgeOptions, nedge.Address(ctx.String("edge_host")))
		}

		if name := ctx.String("edge_transport"); len(name) > 0 && e.opts.EdgeTransport.String() != name {

			if t, ok := e.opts.Edge.Options().Transports[name]; ok {
				e.opts.EdgeTransport = t()
				edgeOptions = append(edgeOptions, nedge.Transport(e.opts.EdgeTransport))
			}
		}

		if e.opts.Action != nil {
			e.opts.Action(ctx)
		}

		return nil
	}))

	e.opts.Service.Init(serviceOpts...)

	e.opts.Edge.Init(edgeOptions...)

	return nil
}

//Run to launch edge server process
func (e *edgeApp) Run() error {

	// init edge itself
	log.Init(log.WithFields(map[string]interface{}{"service": "edge"}))

	// Init plugins
	// for _, p := range Plugins() {
	// 	p.Init(ctx)
	// }

	if err := e.opts.Edge.Run(); err != nil {
		log.Fatal(err)
	}

	// Run go-micro servier
	if err := e.opts.Service.Run(); err != nil {
		log.Fatal(err)
	}

	if err := e.opts.Edge.Stop(); err != nil {
		log.Fatal(err)
	}

	return nil
}

// return micro.service
func (e *edgeApp) MicroService() micro.Service {
	return e.opts.Service
}

// return micro.service
func (e *edgeApp) EdgeService() nedge.Service {
	return e.opts.Edge
}

func (e *edgeApp) Name() string {
	return e.opts.Name
}

func (e *edgeApp) String() string {
	return "edgeApp"
}
