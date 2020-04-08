package main

import (
	//we should use this not x-edge/edge which is a internal service
	edge "github.com/micro-community/x-edge"
	//	_ "github.com/micro-community/x-edge/subscriber"

	log "github.com/micro/go-micro/v2/logger"
)

//XEDGEADDR for target edge address
const XEDGEADDR = "XMicroEdgeServiceAddr"

//XEDGETRANSPORT for target edge port
const XEDGETRANSPORT = "XMicroEdgeServiceTransport"

//Meta Data
var (
	Name    = "go.micro.edge"
	Address = ":8080"
)

func main() {
	// Register Handler
	//protocol.RegisterProtocolHandler(service.Server(), new(handler.Protocol))
	// Register Subscriber
	//eventbroker.RegisterMessageSubscriber(service)

	// Register Publisher
	//eventbroker.RegisterMessagePublisher(service)

	srv := edge.NewService()

	srv.Init()

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
