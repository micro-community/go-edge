package udp

import (
	"bufio"
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
	//寻找确定disconnected的错误，t.conn代表一个实际的连接
	//替代NEWScanner的错误
	//scanner disconnected的错误
	scanner := bufio.NewScanner(u.conn)

	scanner.Split(u.dataExtractor)

	if scanner.Scan() {
		m.Body = scanner.Bytes()
		return nil
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

	return u.conn.Close()

}
