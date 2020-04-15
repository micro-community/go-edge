package cmd

import (
	"github.com/micro-community/x-edge/config"
	ccli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/micro/v2/plugin"
	"github.com/micro/micro/v2/plugin/build"
)

func init() {
	initApp()
}

// initApp initialize edge app to do some basic config
func initApp(options ...micro.Option) {

	setupApp(cmd.App(), options...)

	cmd.Init(
		cmd.Name(config.Name),
		cmd.Description(config.Description),
		cmd.Version(config.BuildVersion()),
	)
}

// setupEdgeApp a edge App to boost
func setupApp(app *ccli.App, options ...micro.Option) {

	// Add the various commands,
	// only plugin(middle ware enabled) for edge
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
			Name:    "edge_web_address",
			Usage:   "Set the edge UI address e.g 0.0.0.0:8082",
			EnvVars: []string{"EDGE_WEB_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "edge_host",
			Usage:   "Set the edge host e.g localhost:8000",
			EnvVars: []string{"EDGE_HOST"},
		},
		&ccli.StringFlag{
			Name:    "edge_namespace",
			Usage:   "Set the namespace used by the edge proxy e.g. com.iot.edge",
			EnvVars: []string{"EDGE_WEB_NAMESPACE"},
		},
		&ccli.StringFlag{
			Name:    "edge_url",
			Usage:   "Set the host used for the edge dashboard e.g edge.project.com",
			EnvVars: []string{"EDGE_WEB_HOST"},
		},

		&ccli.StringFlag{
			Name:    "edge_root_namespace",
			Usage:   "Set the edge root service namespace",
			EnvVars: []string{"EDGE_ROOT_NAMESPACE"},
			Value:   "edge.root",
		},
		&ccli.StringFlag{
			Name:    "edge_transport",
			Usage:   "Set the edge transport to use ,only tcp or udp for now",
			EnvVars: []string{"EDGE_TRANSPORT"},
			Value:   "udp",
		},
	)

	plugins := plugin.Plugins()

	for _, p := range plugins {
		if flags := p.Flags(); len(flags) > 0 {
			app.Flags = append(app.Flags, flags...)
		}
	}

	before := app.Before

	app.Before = func(ctx *ccli.Context) error {

		for _, p := range plugins {
			if err := p.Init(ctx); err != nil {
				return err
			}
		}
		// now do previous before
		return before(ctx)
	}
}
