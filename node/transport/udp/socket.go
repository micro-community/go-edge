package udp

import (
	"bufio"
	"errors"
	"net"
	"time"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/transport"
)

func (u *udpSocket) Local() string {
	return u.conn.LocalAddr().String()
}

func (u *udpSocket) Remote() string {
	return u.dstAddr.String()
	//return u.conn.RemoteAddr().String()
}

func (u *udpSocket) Recv(m *transport.Message) error {
	if m == nil {
		return errors.New("message passed in is nil")
	}
	// set timeout if its greater than 0
	if u.timeout > time.Duration(0) {
		u.conn.SetDeadline(time.Now().Add(u.timeout))
	}

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

	if len(u.packageBuf) > 0 {
		m.Body = u.packageBuf
		u.packageBuf = nil
		//u.packageBuf = u.packageBuf[:0]
		u.packageLen = 0
		//
	} else {
		u.closed = true
		return errors.New("Udp Recv buf is empty")
		//return nil
	}

	return nil
}

func (u *udpSocket) Send(m *transport.Message) error {
	// set timeout if its greater than 0
	if u.timeout > time.Duration(0) {
		u.conn.SetDeadline(time.Now().Add(u.timeout))
	}
	// if err := u.enc.Encode(m); err != nil {
	// 	return err
	// }

	writer := bufio.NewWriter(u.conn)
	writer.Write(m.Body)
	return writer.Flush()

	//return u.encBuf.Flush()
}

func (u *udpSocket) Close() error {
	if u.closed == true {
		u.closed = false
	} else {
		u.conn.Close()
	}

	return nil
}
