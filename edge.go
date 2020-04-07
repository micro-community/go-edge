package edge

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/auth"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/service"
)

//basic metadata
var (
	Name         = "x-edge-service"
	Address      = ":8080"
	Handler      = "meta"
	Resolver     = "micro"
	Namespace    = "x.edge"
	HeaderPrefix = "x-edge-"
)

//Service of edge srv
type Service interface {
	Name() string
	Init(opts ...Option) error
	Run(ctx *cli.Context, srvOpts ...service.Option) error
	String() string
}

//edgeApp for edge process
type edgeApp struct {

	// auth service
	auth auth.Auth
}

//Run to launch edge server process
func (e *edgeApp) Run(ctx *cli.Context, srvOpts ...service.Option) error {
	log.Init(log.WithFields(map[string]interface{}{"service": "edge srv"}))
	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}
	return nil
}
