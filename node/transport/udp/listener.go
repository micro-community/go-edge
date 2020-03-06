package udp

import (
	"bufio"
	"net"
	"time"

	log "github.com/micro/go-micro/v2/logger"
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

		buf := make([]byte, defaultUDPMaxPackageLenth)
		//conn, err := u.listener.Accept()
		bytesLenth, fromAddr, err := u.listener.ReadFromUDP(buf)
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
				log.Infof("udp: Accept error: %v; retrying in %v\n", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		encBuf := bufio.NewWriter(u.listener)
		sock := &udpSocket{
			timeout: u.timeout,
			conn:    u.listener,
			encBuf:  encBuf,
			//enc:        gob.NewEncoder(encBuf),
			//dec:        gob.NewDecoder(u.listener),
			packageBuf: buf,
			packageLen: bytesLenth,
			dstAddr:    fromAddr,
			closed:     false,
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
