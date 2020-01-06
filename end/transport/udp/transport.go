package udp

import (
	"bufio"
	"encoding/gob"
	"net"

	"github.com/micro/go-micro/transport"
)

func (u *udpTransport) Dial(addr string, opts ...transport.DialOption) (transport.Client, error) {
	dopts := transport.DialOptions{
		Timeout: transport.DefaultDialTimeout,
	}

	for _, opt := range opts {
		opt(&dopts)
	}

	var conn net.Conn

	// ctx, _ := context.WithTimeout(context.Background(), dopts.Timeout)
	// conn, err = u.dialContext(ctx, addr)

	conn, err := net.DialTimeout("udp", addr, dopts.Timeout)

	if err != nil {
		return nil, nil
	}

	encBuf := bufio.NewWriter(conn)

	return &udpClient{
		dialOpts: dopts,
		conn:     conn,
		encBuf:   encBuf,
		enc:      gob.NewEncoder(encBuf),
		dec:      gob.NewDecoder(conn),
		timeout:  u.opts.Timeout,
	}, nil
}

func (u *udpTransport) Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	var options transport.ListenOptions
	for _, o := range opts {
		o(&options)
	}
	//	var l net.Listener
	var err error

	udpAdress, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenUDP("udp", udpAdress)

	if err != nil {
		return nil, err
	}

	return &udpListener{
		timeout:  u.opts.Timeout,
		listener: l,
		opts:     options,
	}, nil
}

func (u *udpTransport) Init(opts ...transport.Option) error {
	for _, o := range opts {
		o(&u.opts)
	}
	return nil
}

func (u *udpTransport) Options() transport.Options {
	return u.opts
}

func (u *udpTransport) String() string {
	return "udp"
}
