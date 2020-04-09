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
	Name         = "x-edge-app"
	Address      = ":8000"
	Transport    = "udp"
	Host         = ":8080"
	Namespace    = "x.edge"
	HeaderPrefix = "x-edge-"
)

//Service of edge srv
type Service interface {
	Name() string
	Init(opts ...Option) error
	Run() error
	String() string
}

//edgeApp for edge process
type edgeApp struct {
	opts Options
	// auth service
	auth auth.Auth
}

//NewService return a edge service application
func NewService() Service {

	return nil
}

func (e *edgeApp) buildGoMicro() []micro.Option {

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

		if e.opts.Action != nil {
			e.opts.Action(ctx)
		}

		return nil
	}))

	return serviceOpts

}

func (e *edgeApp) Init(opts ...Option) error {
	for _, o := range opts {
		o(&e.opts)
	}

	serviceOpts := e.buildGoMicro()
	microOptions := []nedge.Option{}

	edgeOptions := append(microOptions, micro.Action(func(ctx *cli.Context) error {

		if len(ctx.String("edge_transport")) > 0 {
			Name = ctx.String("edge_transport")
		}
		if len(ctx.String("edge_web_address")) > 0 {
			Address = ctx.String("edge_web_address")
		}
		if len(ctx.String("edge_namespace")) > 0 {
			Namespace = ctx.String("namespace")
		}

		if e.opts.Action != nil {
			e.opts.Action(ctx)
		}
		return nil
	}))
	e.opts.Edge.Init()
	e.opts.Service.Init(edgeOptions...)
	e.opts.Service.Init(serviceOpts...)
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

	// pass namespace and resolver through to the server as these are needed to perform auth
	edgesrv := nedge.NewServer(nedge.Address(Address), nedge.Namespace(Namespace))

	if err := edgesrv.Run(); err != nil {
		log.Fatal(err)
	}

	// Run go-micro servier
	if err := e.opts.Service.Run(); err != nil {
		log.Fatal(err)
	}

	if err := edgesrv.Stop(); err != nil {
		log.Fatal(err)
	}

	return nil
}

// return micro.service
func (e *edgeApp) MicroService() micro.Service {
	return e.Options().Service
}
