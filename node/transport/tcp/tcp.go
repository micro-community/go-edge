//Package tcp provides a TCP transport
package tcp

import (
	"errors"

	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/transport"
)

var (
	errorTransportDataExtract = errors.New("extract data error in tcp transport")
)


func init() {
	cmd.DefaultTransports["tcp"] = NewTransport
}


//NewTransport Return a New TCP Transport
func NewTransport(opts ...transport.Option) transport.Transport {
	var options transport.Options
	for _, o := range opts {
		o(&options)
	}
	return &tcpTransport{opts: options}
}
