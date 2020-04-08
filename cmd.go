package edge

import (
	"github.com/micro-community/x-edge/config"
	ccli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/micro/v2/plugin"
	"github.com/micro/micro/v2/plugin/build"
)

func init() {
	// // setup the build plugin
	plugin.Register(build.Flags())
}

// Init initialised the command line
func Init(options ...micro.Option) {
	Setup(cmd.App(), options...)

	cmd.Init(
		cmd.Name(config.Name),
		cmd.Description(config.Description),
		cmd.Version(config.BuildVersion()),
	)
}

// Setup a cli.App
func Setup(app *ccli.App, options ...micro.Option) {

	// Add the various commands,
	// only plugin for edge
	app.Commands = append(app.Commands, build.Commands()...)

	setup(app)
}

//setup edge command lines
func setup(app *ccli.App) {

	app.Flags = append(app.Flags,
		&ccli.BoolFlag{
			Name:  "local",
			Usage: "Enable local only development: Defaults to true.",
		},
		&ccli.StringFlag{
			Name:    "edge_address",
			Usage:   "Set the edge UI address e.g 0.0.0.0:8082",
			EnvVars: []string{"EDGE_WEB_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "edge_namespace",
			Usage:   "Set the namespace used by the edge proxy e.g. hw.hbt.edge",
			EnvVars: []string{"EDGE_WEB_NAMESPACE"},
		},
		&ccli.StringFlag{
			Name:    "edge_url",
			Usage:   "Set the host used for the edge dashboard e.g edge.example.com",
			EnvVars: []string{"EDGE_WEB_HOST"},
		},

		&ccli.StringFlag{
			Name:    "edge_root_namespace",
			Usage:   "Set the edge root service namespace",
			EnvVars: []string{"EDGE_ROOT_NAMESPACE"},
			Value:   "edge.root",
		},
	)

	plugins := plugin.Plugins()

	for _, p := range plugins {
		if flags := p.Flags(); len(flags) > 0 {
			app.Flags = append(app.Flags, flags...)
		}

		if cmds := p.Commands(); len(cmds) > 0 {
			app.Commands = append(app.Commands, cmds...)
		}
	}

	before := app.Before

	app.Before = func(ctx *ccli.Context) error {

		if len(ctx.String("edge_address")) > 0 {
			//	edge.Address = ctx.String("edge_address")
		}
		if len(ctx.String("edge_namespace")) > 0 {
			//	edge.Namespace = ctx.String("edge_namespace")
		}
		if len(ctx.String("edge_host")) > 0 {
			//	edge.Host = ctx.String("edge_host")
		}

		for _, p := range plugins {
			if err := p.Init(ctx); err != nil {
				return err
			}
		}

		// now do previous before
		return before(ctx)
	}
}
