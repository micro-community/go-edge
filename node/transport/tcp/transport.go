
//Package tcp provides a TCP transport
package tcp

import (
	"bufio"
	"crypto/tls"
	"net"
	"github.com/micro/go-micro/transport"
	maddr "github.com/micro/go-micro/util/addr"
	mnet "github.com/micro/go-micro/util/net"
	mls "github.com/micro/go-micro/util/tls"
)


type tcpTransport struct {
	opts transport.Options
}


func (t *tcpTransport) Dial(addr string, opts ...transport.DialOption) (transport.Client, error) {
	dopts := transport.DialOptions{
		Timeout: transport.DefaultDialTimeout,
	}

	for _, opt := range opts {
		opt(&dopts)
	}

	var conn net.Conn
	var err error

	// TODO: support dial option here rather than using internal config
	if t.opts.Secure || t.opts.TLSConfig != nil {
		config := t.opts.TLSConfig
		if config == nil {
			config = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: dopts.Timeout}, "tcp", addr, config)
	} else {
		conn, err = net.DialTimeout("tcp", addr, dopts.Timeout)
	}

	if err != nil {
		return nil, err
	}

	encBuf := bufio.NewWriter(conn)

	return &tcpTransportClient{
		dialOpts: dopts,
		conn:     conn,
		encBuf:   encBuf,
		//		enc:      gob.NewEncoder(encBuf),
		//		dec:      gob.NewDecoder(conn),
		timeout: t.opts.Timeout,
	}, nil
}

func (t *tcpTransport) Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	var options transport.ListenOptions
	for _, o := range opts {
		o(&options)
	}

	var l net.Listener
	var err error

	// TODO: support use of listen options
	if t.opts.Secure || t.opts.TLSConfig != nil {
		config := t.opts.TLSConfig

		fn := func(addr string) (net.Listener, error) {
			if config == nil {
				hosts := []string{addr}

				// check if its a valid host:port
				if host, _, err := net.SplitHostPort(addr); err == nil {
					if len(host) == 0 {
						hosts = maddr.IPs()
					} else {
						hosts = []string{host}
					}
				}

				// generate a certificate
				cert, err := mls.Certificate(hosts...)
				if err != nil {
					return nil, err
				}
				config = &tls.Config{Certificates: []tls.Certificate{cert}}
			}
			return tls.Listen("tcp", addr, config)
		}

		l, err = mnet.Listen(addr, fn)
	} else {
		fn := func(addr string) (net.Listener, error) {
			return net.Listen("tcp", addr)
		}

		l, err = mnet.Listen(addr, fn)
	}

	if err != nil {
		return nil, err
	}

	return &tcpTransportListener{
		timeout:  t.opts.Timeout,
		listener: l,
	}, nil
}

func (t *tcpTransport) Init(opts ...transport.Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *tcpTransport) Options() transport.Options {
	return t.opts
}

func (t *tcpTransport) String() string {
	return "tcp"
}