package edge

import(
	"github.com/micro/go-micro/v2/transport"
	"github.com/micro-community/x-edge/node/transport/udp"
	"github.com/micro-community/x-edge/node/transport/tcp"
)

var(

	DefaultTransports = map[string]func(...transport.Option) transport.Transport{
		"udp": udp.NewTransport,
		"tcp":   tcp.NewTransport,
	}
)