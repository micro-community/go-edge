package udp

import (
	"errors"
	"time"

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

	if len(u.packageBuf) > 0 {
		m.Body = u.packageBuf
		u.packageBuf = nil
		u.packageLen = 0
	} else {
		u.closed = true
		return errors.New("Udp Recv buf is empty")
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
	u.conn.WriteTo(m.Body, u.dstAddr)
	//writer := bufio.NewWriter(u.conn)
	//writer.Write(m.Body)
	//return writer.Flush()
	return nil
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
