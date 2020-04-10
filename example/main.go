package main

import (
	//we should use this not x-edge/edge which is a internal service
	edge "github.com/micro-community/x-edge"
	"github.com/micro-community/x-edge/node/transport/udp"
	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
)

//Meta Data
var (
	Name    = "go.micro.x.edge.example"
	Address = ":8080"
)

func main() {
	// Register Handler
	//protocol.RegisterProtocolHandler(service.Server(), new(handler.Protocol))
	// Register Subscriber
	//eventbroker.RegisterMessageSubscriber(service)
	// Register Publisher
	//eventbroker.RegisterMessagePublisher(service)

	srv := edge.NewService(edge.EgTransport(udp.NewTransport()), edge.Version("v1.0.0"))

	srv.Init(edge.Action(func(ctx *cli.Context) {

		// here, do your own
		if name := ctx.String("server_name"); len(name) > 0 {
			Name = name
		}

	}))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
