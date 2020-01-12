package main

import (
	"os"
	"regexp"

	edge "github.com/micro-community/x-edge"
	eventbroker "github.com/micro-community/x-edge/broker"
	"github.com/micro-community/x-edge/config"
	"github.com/micro-community/x-edge/end/transport/extractor"
	"github.com/micro-community/x-edge/handler"

	protocol "github.com/micro-community/x-edge/proto/protocol"
	_ "github.com/micro-community/x-edge/subscriber"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
)

//XEDGEADDR for target edge address
const XEDGEADDR = "XMicroEdgeServiceAddr"

//XEDGETRANSPORT for target edge port
const XEDGETRANSPORT = "XMicroEdgeServiceTransport"

func main() {
	// New Service
	service := micro.NewService(
		//Select transport protocol (eg:tcp or udp) for XMicroEdgeService
		micro.Flags(cli.StringFlag{
			Name:   "XMicroEdgeServiceTransport",
			Usage:  "tcp",
			EnvVar: XEDGETRANSPORT,
			Value:  "tcp", //or udp
		}),
		//Set address for XMicroEdgeService 192.168.1.198:6600
		micro.Flags(cli.StringFlag{
			Name:   "XMicroEdgeServiceAddr",
			Usage:  "format: 127.0.0.1:6600",
			EnvVar: XEDGEADDR,
			//Value:  "192.168.1.198:6600",
			Value: "192.168.1.198:1234",
		}),
	)

	// Initialise service
	service.Init(
		micro.Action(func(c *cli.Context) {
			if info := c.String("XMicroEdgeServiceTransport"); info != "" {
				log.Log("XMicroEdgeServiceTransport:", info)
				config.XMicroEdgeServiceTransport = info
			} else {
				if env := os.Getenv(XEDGETRANSPORT); len(env) > 0 {
					log.Log(XEDGETRANSPORT, ":", env)
					config.XMicroEdgeServiceTransport = env
				} else {
					log.Log("default XMicroEdgeServiceTransport is tcp")
				}
			}

			if info := c.String("XMicroEdgeServiceAddr"); info != "" {
				log.Log("XMicroEdgeServiceAddr:", info)
				config.XMicroEdgeServiceAddr = info
			} else {
				if env := os.Getenv(XEDGEADDR); len(env) > 0 {
					log.Log(XEDGEADDR, ":", env)
					config.XMicroEdgeServiceAddr = env
				} else {
					log.Log("default XMicroEdgeServiceAddr is 192.168.1.198:6600")
				}
			}
		}),
	)
	// Register Handler for Data Extractor
	extractor.RegisterExtractorHandler(DataExtractor)

	// Register Handler
	protocol.RegisterProtocolHandler(service.Server(), new(handler.Protocol))

	// Register Subscriber
	//eventbroker.RegisterMessageSubscriber(service)

	// Register Publisher
	eventbroker.RegisterMessagePublisher(service)

	//run the second listening service， you could set the args in config
	//config.XMicroEdgeServiceTransport
	//config.XMicroEdgeServiceAddr

	go edge.RunProc()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func DataExtractor(data []byte, atEOF bool) (advance int, token []byte, err error) {
	//set package min lengh
	minDataPackageLenth := 50

	if atEOF || len(data) == 0 {
		return 0, nil, nil
	}

	reg, _ := regexp.Compile("(?i:</protocol>)")

	indexs := reg.FindIndex(data)

	if indexs == nil || indexs[0] <= minDataPackageLenth {
		return -1, data, nil //errors.New("error to extract data from socket")
	}

	advance = indexs[1]
	token = data[0:indexs[1]]
	return
}
