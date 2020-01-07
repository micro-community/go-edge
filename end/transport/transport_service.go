package transport

import (
	"strings"
	tcp "github.com/micro-community/x-edge/end/transport/tcp"
	udp "github.com/micro-community/x-edge/end/transport/udp"
	ts "github.com/micro/go-micro/transport"
)

func CreateTransport(name string) ts.Transport {
	str := strings.ToLower(name)
	if str == "udp" {
		return udp.NewTransport()
	} else {
		return tcp.NewTransport()
	}

}
