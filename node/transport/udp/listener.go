package udp

import (
	"time"

	"github.com/micro/go-micro/v2/transport"
)

//　　UDP : 1500 - IP(20) - UDP(8) = 1472(Bytes)
const defaultUDPMaxPackageLenth = 1472

func (u *udpListener) Addr() string {
	return u.listener.LocalAddr().String()
}

func (u *udpListener) Close() error {
	return u.listener.Close()
}

//Accept and handle a data package
func (u *udpListener) Accept(fn func(transport.Socket)) error {
	var tempDelay time.Duration

	for {

		select {
		case <-u.sockexit:
			return nil
			//			encBuf := bufio.NewWriter(u.listener)
		case c := <-u.conn:
			sock := &udpSocket{
				timeout: u.opts.Timeout,
				ctx:     u.opts.Context,
				conn:    u.listener,
				local:   c.Remote(),
				remote:  c.Local(),
				closed:  c.exit,
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
}
