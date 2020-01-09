package edge

import (
	"sync"

	"github.com/micro-community/x-edge/config"
	esrv "github.com/micro-community/x-edge/node/server"
	"github.com/micro-community/x-edge/node/transport/tcp"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
)

type service struct {
	opts Options
	//	mux *http.ServeMux
	sync.Mutex
	running bool
	static  bool
	exit    chan chan error
}

func newService(opts ...Option) micro.Service {
	options := newOptions(opts...)

	svr := esrv.NewServer(server.WithRouter(esrv.DefaultRouter()))

	nodeService := micro.NewService(
		micro.Server(svr),
		micro.Name(config.XMicroEdgeServiceName),
		micro.Address(config.XMicroEdgeServiceAddr),
		micro.Transport(tcp.NewTransport()),
		micro.Metadata(map[string]string{
			"type": "protocol-edge-node-server",
		},
		),
	)

	options.Service = nodeService

	s := &service{
		opts:   options,
		static: true,
	}

	return s
}

// //RunProc Listen to The Device Edge Server (Gateway、Controller etc)
// func RunProc() {

// 	r := esrv.DefaultRouter()

// 	svr := esrv.NewServer(server.WithRouter(r))

// 	t := ts.CreateTransport(config.XMicroEdgeServiceTransport)
// 	// Create a new service. Optionally include some options here.
// 	service := micro.NewService(
// 		micro.Server(svr),
// 		micro.Name(config.XMicroEdgeServiceName),
// 		micro.Version("latest"),
// 		micro.Address(config.XMicroEdgeServiceAddr),
// 		micro.Transport(t),
// 		micro.Metadata(map[string]string{
// 			"type": "protocol-controller-server",
// 		},
// 		),
// 	)

// 	service.Init()

// 	//register server message router
// 	router.RegisterProtocolHandler(svr, new(router.ProtocolServer))

// 	// Run the server
// 	if err := service.Run(); err != nil {
// 		fmt.Println(err)
// 	}
// }
