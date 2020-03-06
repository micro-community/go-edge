package main

import (
	"os"

	edge "github.com/micro-community/x-edge"
	eventbroker "github.com/micro-community/x-edge/broker"
	"github.com/micro-community/x-edge/config"
	"github.com/micro-community/x-edge/handler"
	protocol "github.com/micro-community/x-edge/proto/protocol"
	_ "github.com/micro-community/x-edge/subscriber"
	cli "github.com/micro/cli/v2"
	micro "github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)

//XEDGEADDR for target edge address
const XEDGEADDR = "XMicroEdgeServiceAddr"

//XEDGETRANSPORT for target edge port
const XEDGETRANSPORT = "XMicroEdgeServiceTransport"

func main() {
	// New Service
	service := micro.NewService(
		//Select transport protocol (eg:tcp or udp) for XMicroEdgeService
		micro.Flags(
			&cli.StringFlag{
				Name:    "XMicroEdgeServiceTransport",
				Usage:   "tcp",
				EnvVars: []string{XEDGETRANSPORT},
				//Value: "tcp"
				Value: "ppp",
			},
			&cli.StringFlag{
				Name:    "XMicroEdgeServiceAddr",
				Usage:   "format: 127.0.0.1:6600",
				EnvVars: []string{XEDGEADDR},
				//Value:  "192.168.1.198:6600",
				Value: "192.168.1.198:1234",
			},
		),
	)

	// Initialise service
	service.Init(
		micro.Action(func(c *cli.Context) error {
			if info := c.String("XMicroEdgeServiceTransport"); info != "" {
				log.Info("XMicroEdgeServiceTransport:", info)
				config.XMicroEdgeServiceTransport = info
			} else {
				if env := os.Getenv(XEDGETRANSPORT); len(env) > 0 {
					log.Info(XEDGETRANSPORT, ":", env)
					config.XMicroEdgeServiceTransport = env
				} else {
					log.Info("default XMicroEdgeServiceTransport is tcp")
				}
			}

			if info := c.String("XMicroEdgeServiceAddr"); info != "" {
				log.Info("XMicroEdgeServiceAddr:", info)
				config.XMicroEdgeServiceAddr = info
			} else {
				if env := os.Getenv(XEDGEADDR); len(env) > 0 {
					log.Info(XEDGEADDR, ":", env)
					config.XMicroEdgeServiceAddr = env
				} else {
					log.Info("default XMicroEdgeServiceAddr is 192.168.1.198:6600")
				}
			}
			return nil
		}),
	)

	// Register Handler
	protocol.RegisterProtocolHandler(service.Server(), new(handler.Protocol))

	// Register Subscriber
	//eventbroker.RegisterMessageSubscriber(service)

	// Register Publisher
	eventbroker.RegisterMessagePublisher(service)

	//run the second listening service， you could set the args in config
	//config.XMicroEdgeServiceTransport
	//config.XMicroEdgeServiceAddr

	srv := edge.NewService()
	srv.Run()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
