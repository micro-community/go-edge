package udp

import (
	"time"

	"github.com/micro/go-micro/transport"
)

func (u *udpClient) Local() string {
	return u.conn.LocalAddr().String()
}

func (u *udpClient) Remote() string {
	return u.conn.RemoteAddr().String()
}

func (u *udpClient) Send(m *transport.Message) error {
	// set timeout if its greater than 0
	if u.timeout > time.Duration(0) {
		u.conn.SetDeadline(time.Now().Add(u.timeout))
	}
	if err := u.enc.Encode(m); err != nil {
		return err
	}
	return u.encBuf.Flush()
}

func (u *udpClient) Recv(m *transport.Message) error {
	// set timeout if its greater than 0
	if u.timeout > time.Duration(0) {
		u.conn.SetDeadline(time.Now().Add(u.timeout))
	}
	return u.dec.Decode(&m)
}

func (u *udpClient) Close() error {
	return u.conn.Close()
}
