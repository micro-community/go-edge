package controller

import (
	"fmt"
	esrv "github.com/micro-community/x-micro-edge/end/server"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"

	"github.com/micro-community/x-micro-edge/config"
	"github.com/micro-community/x-micro-edge/end/router"
	ts "github.com/micro-community/x-micro-edge/end/transport"
)

const XEDGEADDR = "XMicroEdgeServiceAddr"
const XEDGETRANSPORT = "XMicroEdgeServiceTransport"

//RunProc Listen the Controller
func RunProc() {

	r := esrv.DefaultRouter()

	svr := esrv.NewServer(server.WithRouter(r))

	t := ts.CreateTransport(config.XMicroEdgeServiceTransport)
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Server(svr),
		micro.Name(config.XMicroEdgeServiceName),
		micro.Version("latest"),
		micro.Address(config.XMicroEdgeServiceAddr),
		micro.Transport(t),
		micro.Metadata(map[string]string{
			"type": "protocol-controller-server",
		},
		),
	)

	service.Init()

	//register server message router
	router.RegisterProtocolHandler(svr, new(router.ProtocolServer))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
