//Package tcp provides a TCP transport
package tcp

import (
	"bufio"
	"errors"
	"net"
	"time"

	nts "github.com/micro-community/x-edge/node/transport"
	"github.com/micro/go-micro/transport"
)

type tcpTransportClient struct {
	dialOpts transport.DialOptions
	conn     net.Conn
	//	enc      *gob.Encoder
	//	dec      *gob.Decoder
	encBuf        *bufio.Writer
	timeout       time.Duration
	dataExtractor nts.DataExtractor
}

func (t *tcpTransportClient) Local() string {
	return t.conn.LocalAddr().String()
}

func (t *tcpTransportClient) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *tcpTransportClient) Send(m *transport.Message) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}
	writer := bufio.NewWriter(t.conn)
	writer.Write(m.Body)
	return writer.Flush()
}

func (t *tcpTransportClient) Recv(m *transport.Message) error {
	// set timeout if its greater than 0

	if t.dialOpts.Context != nil && t.dataExtractor == nil {
		if v := t.dialOpts.Context.Value(nts.DataExtractorFuncKey{}); v != nil {
			t.dataExtractor = v.(nts.DataExtractor)
		}
	}

	if m == nil {
		return errors.New("message passed in is nil")
	}

	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}

	scanner := bufio.NewScanner(t.conn)
	scanner.Split(t.dataExtractor)

	if scanner.Scan() {
		m.Body = scanner.Bytes()
		return nil
	}

	return errorTransportDataExtract

}

func (t *tcpTransportClient) Close() error {
	return t.conn.Close()
}