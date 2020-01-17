//Package tcp provides a TCP transport
package tcp

import (
	"bufio"
	"crypto/tls"
	"errors"
	"net"
	"time"
	nts "github.com/micro-community/x-edge/node/transpot"
	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/transport"
	maddr "github.com/micro/go-micro/util/addr"
	"github.com/micro/go-micro/util/log"
	mnet "github.com/micro/go-micro/util/net"
	mls "github.com/micro/go-micro/util/tls"
)

var (
	errorTransportDataExtract = errors.New("extract data error in tcp transport")
)

type tcpTransport struct {
	opts transport.Options
}

type tcpTransportClient struct {
	dialOpts transport.DialOptions
	conn     net.Conn
	//	enc      *gob.Encoder
	//	dec      *gob.Decoder
	encBuf  *bufio.Writer
	timeout time.Duration
	dataExtractor nts.DataExtractor
}

type tcpTransportSocket struct {
	conn net.Conn
	//	enc     *gob.Encoder
	//	dec     *gob.Decoder
	encBuf  *bufio.Writer
	timeout time.Duration
}

type tcpTransportListener struct {
	listener net.Listener
	timeout  time.Duration
}

func init() {
	cmd.DefaultTransports["tcp"] = NewTransport
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
		if v := t.dialOpts.Context.Value(codecsKey{}); v != nil {
			t.dataExtractor = v.(transport.DataExtractor)
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

	scanner.Split(dataExtractor)

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

func (t *tcpTransportListener) Addr() string {
	return t.listener.Addr().String()
}

func (t *tcpTransportListener) Close() error {
	return t.listener.Close()
}

func (t *tcpTransportListener) Accept(fn func(transport.Socket)) error {
	var tempDelay time.Duration

	for {
		c, err := t.listener.Accept()
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
				log.Logf("http: Accept error: %v; retrying in %v\n", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		encBuf := bufio.NewWriter(c)
		sock := &tcpTransportSocket{
			timeout: t.timeout,
			conn:    c,
			encBuf:  encBuf,
			//			enc:     gob.NewEncoder(encBuf),
			//			dec:     gob.NewDecoder(c),
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


//NewTransport Return a New TCP Transport
func NewTransport(opts ...transport.Option) transport.Transport {
	var options transport.Options
	for _, o := range opts {
		o(&options)
	}
	return &tcpTransport{opts: options}
}
