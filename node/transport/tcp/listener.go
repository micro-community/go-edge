//Package tcp provides a TCP transport
package tcp

import (
	"bufio"
	"net"
	"time"

	nts "github.com/micro-community/x-edge/node/transport"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-micro/util/log"
)

type tcpTransportListener struct {
	listener      net.Listener
	timeout       time.Duration
	dataExtractor nts.DataExtractor
}

func (t *tcpTransportListener) Addr() string {
	return t.listener.Addr().String()
}

func (t *tcpTransportListener) Close() error {
	return t.listener.Close()
}

func (t *tcpTransportListener) Accept(fn func(transport.Socket)) error {
	var tempDelay time.Duration

	for {
		c, err := t.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Logf("http: Accept error: %v; retrying in %v\n", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		encBuf := bufio.NewWriter(c)
		sock := &tcpTransportSocket{
			timeout:       t.timeout,
			conn:          c,
			encBuf:        encBuf,
			dataExtractor: t.dataExtractor,
			//			enc:     gob.NewEncoder(encBuf),
			//			dec:     gob.NewDecoder(c),
		}

		go func() {
			// TODO: think of a better error response strategy
			defer func() {
				if r := recover(); r != nil {
					sock.Close()
				}
			}()

			fn(sock)
		}()
	}
}
