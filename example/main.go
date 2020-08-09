package main

import (
	//we should use this not x-edge/edge which is a internal service
	edge "github.com/micro-community/x-edge"
	"github.com/micro-community/x-edge/node/transport/tcp"
	protocol "github.com/micro-community/x-edge/proto/protocol"
	eventbroker "github.com/micro-community/x-edge/broker"
	handler "github.com/micro-community/x-edge/handler"
	"github.com/micro-community/x-edge/node/router"
	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
)

//Meta Data
var (
	Name    = "go.micro.x.edge.example"
	Address = ":8080"
)

func main() {

	srv := edge.NewService(
		edge.EgTransport(tcp.NewTransport()),
		edge.Version("v1.0.0"),
		edge.EgExtractor(dataExtractor))

	srv.Init(edge.Action(func(ctx *cli.Context) {

		// here, do your own
		if name := ctx.String("server_name"); len(name) > 0 {
			Name = name
		}
	}))

	// Register Handler
	protocol.RegisterProtocolHandler(srv.MService().Server(), new(handler.Protocol))
	//Register Subscriber
	eventbroker.RegisterMessageSubscriber(srv.MService())
	//Register Publisher
	eventbroker.RegisterMessagePublisher(srv.MService())
	//register server message router
	router.RegisterProtocolHandler(srv.EService().Server(), new(router.ProtocolServer))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
