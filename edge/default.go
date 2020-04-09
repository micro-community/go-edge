package edge

import (
	"github.com/micro-community/x-edge/node/transport/tcp"
	"github.com/micro-community/x-edge/node/transport/udp"
	"github.com/micro/go-micro/v2/transport"
)

// default components for edge node server
var (
	DefaultTransports = map[string]func(...transport.Option) transport.Transport{
		"udp": udp.NewTransport,
		"tcp": tcp.NewTransport,
	}
)
