package udp

import (
	"bufio"
	"errors"
	"time"

	"github.com/micro/go-micro/v2/transport"
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
	writer := bufio.NewWriter(u.conn)
	writer.Write(m.Body)
	return writer.Flush()

}

func (u *udpClient) Recv(m *transport.Message) error {
	// set timeout if its greater than 0
	if u.timeout > time.Duration(0) {
		u.conn.SetDeadline(time.Now().Add(u.timeout))
	} else if m == nil {
		return errors.New("message passed in is nil")
	}

	scanner := bufio.NewScanner(u.conn)
	scanner.Split(u.dataExtractor)

	if scanner.Scan() {
		m.Body = scanner.Bytes()
		return nil
	}

	return nil
}

func (u *udpClient) Close() error {
	return u.conn.Close()
}
