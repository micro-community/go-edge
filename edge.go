package edge

import (
	"strings"

	nedge "github.com/micro-community/x-edge/edge"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
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
	MService() micro.Service
	EService() nedge.Service
}

//edgeApp for edge process
type edgeApp struct {
	opts Options
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

	if e.opts.EdgeTransport != nil {
		edgeOptions = append(edgeOptions, nedge.Transport(e.opts.EdgeTransport))
	}
	if strings.Trim(e.opts.EdgeHost, " ") != "" {
		edgeOptions = append(edgeOptions, nedge.Host(e.opts.EdgeHost))
	}
	if e.opts.Extractor != nil {
		edgeOptions = append(edgeOptions, nedge.WithExtractor(e.opts.Extractor))
	}

	e.opts.Edge.Init(edgeOptions...)

	serviceOpts = append(serviceOpts, micro.Action(func(ctx *cli.Context) error {
		// execute edge service Action
		if e.opts.Edge.Options().Action != nil {
			e.opts.Edge.Options().Action(ctx)
		}
		if e.opts.Action != nil {
			e.opts.Action(ctx)
		}

		return nil
	}))

	e.opts.MicroService.Init(serviceOpts...)

	return nil
}

func (e *edgeApp) start() error {

	return nil
}

func (e *edgeApp) stop() error {
	return nil
}

//Run to launch edge server process
func (e *edgeApp) Run() error {

	// init edge itself
	log.Init(log.WithFields(map[string]interface{}{"service": "edge"}))

	if err := e.opts.Edge.Run(); err != nil {
		log.Fatal(err)
	}

	// Run go-micro servier
	if err := e.opts.MicroService.Run(); err != nil {
		log.Fatal(err)
	}

	if err := e.opts.Edge.Stop(); err != nil {
		log.Fatal(err)
	}

	return nil
}

// return micro.service
func (e *edgeApp) MService() micro.Service {
	return e.opts.MicroService
}

// return micro.service
func (e *edgeApp) EService() nedge.Service {
	return e.opts.Edge
}

func (e *edgeApp) Name() string {
	return e.opts.Name
}

func (e *edgeApp) String() string {
	return "edgeApp"
}
