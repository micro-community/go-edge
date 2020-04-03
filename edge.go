package edge

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)

//basic metadata
var (
	Name         = "x-edge-service"
	Address      = ":8080"
	Handler      = "meta"
	Resolver     = "micro"
	Namespace    = "x.edge"
	HeaderPrefix = "X-edge-"
)

//Service of edge srv
type Service interface {
	Name() string
	Init(opts ...Option) error
	Options() Options
	Run(ctx *cli.Context, srvOpts ...micro.Option) error
	String() string
}
type srv struct {
}

//Run to launch edge server
func (s *srv) Run(ctx *cli.Context, srvOpts ...micro.Option) {
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
