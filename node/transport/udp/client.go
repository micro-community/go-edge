package udp

import (
	"errors"
	"io/ioutil"
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
	u.conn.Write(m.Body)

	return nil
	// writer := bufio.NewWriter(u.conn)
	// writer.Write(m.Body)
	// return writer.Flush()

}

func (u *udpClient) Recv(m *transport.Message) error {

	if m == nil {
		return errors.New("message passed in is nil")
	}

	// set timeout if its greater than 0
	if u.timeout > time.Duration(0) {
		u.conn.SetDeadline(time.Now().Add(u.timeout))
	}
	data, err := ioutil.ReadAll(u.conn)
	m.Body = data

	return err
	// scanner := bufio.NewScanner(u.conn)
	// scanner.Split(u.dataExtractor)

	// if scanner.Scan() {
	// 	m.Body = scanner.Bytes()
	// 	return nil
	// }

	//return nil
}

func (u *udpClient) Close() error {
	return u.conn.Close()
}
