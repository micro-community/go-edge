package udp

import (
	"bufio"
	"github.com/micro/go-micro/v2/transport"
	//	"github.com/micro/go-micro/v2/util/ring"
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
	for {
		buf := make([]byte, defaultUDPMaxPackageLenth)
		//	rbuffer := ring.New(defaultUDPMaxPackageLenth),
		//conn, err := u.listener.Accept()
		//bytesLenth, fromAddr, err := u.listener.ReadFromUDP(buf)
		bytesLennth, fromAddr, err := u.pConn.ReadFrom(buf)

		if err != nil {
			u.errorChan <- err
		}
		select {
		case <-u.exit:
			return nil
		case <-u.errorChan:
			return nil
		default:
			sock := &udpSocket{
				timeout: u.opts.Timeout,
				conn:    u.listener,
				pConn:   u.listener,
				remote:  fromAddr.String(),
				local:   u.Addr(),
				encBuf: bufio.NewWriter(u.listener)
				exit:   make(chan bool)
			}
			go fn(sock)

		}
	}
}
