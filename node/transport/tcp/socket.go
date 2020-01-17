//Package tcp provides a TCP transport
package tcp

import (
	"bufio"
	"errors"
	nts "github.com/micro-community/x-edge/node/transport"
	"github.com/micro/go-micro/transport"
	"net"
	"time"
)

type tcpTransportSocket struct {
	conn net.Conn
	//	enc     *gob.Encoder
	//	dec     *gob.Decoder
	encBuf        *bufio.Writer
	timeout       time.Duration
	dataExtractor nts.DataExtractor
}

func (t *tcpTransportSocket) Local() string {
	return t.conn.LocalAddr().String()
}

func (t *tcpTransportSocket) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *tcpTransportSocket) Recv(m *transport.Message) error {
	if m == nil {
		return errors.New("message passed in is nil")
	}
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}
	//寻找确定disconnected的错误，t.conn代表一个实际的连接
	//替代NEWScanner的错误
	//scanner disconnected的错误
	scanner := bufio.NewScanner(t.conn)

	scanner.Split(t.dataExtractor)

	if scanner.Scan() {
		m.Body = scanner.Bytes()
		return nil
	}

	return errorTransportDataExtract
}

func (t *tcpTransportSocket) Send(m *transport.Message) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}

	writer := bufio.NewWriter(t.conn)
	writer.Write(m.Body)
	return writer.Flush()

	//_, err := t.conn.Write(m.Body)
	//return err
}

func (t *tcpTransportSocket) Close() error {
	return t.conn.Close()
}
