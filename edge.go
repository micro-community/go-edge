package edge

import (
	nedge "github.com/micro-community/x-edge/edge"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/auth"
	log "github.com/micro/go-micro/v2/logger"
)

//basic metadata
var (
	Name         = "x-edge-app"
	Address      = ":8000"
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
func (e *edgeApp) Init(opts ...Option) error {

	return nil
}

//Run to launch edge server process
func (e *edgeApp) Run(ctx *cli.Context) error {

	log.Init(log.WithFields(map[string]interface{}{"service": "edge"}))
	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("namespace")) > 0 {
		Namespace = ctx.String("namespace")
	}
	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// pass namespace and resolver through to the server as these are needed to perform auth
	edgesrv := nedge.NewServer(nedge.Address(Address), nedge.Namespace(Namespace))

	//get internal service
	service := e.opts.Service

	//setup some thing
	service.Init()

	if err := edgesrv.Run(); err != nil {
		log.Fatal(err)
	}

	// Run go-micro servier
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	if err := edgesrv.Stop(); err != nil {
		log.Fatal(err)
	}

	return nil
}
