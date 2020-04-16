package udp

import (
	"bufio"
	"net"

	"github.com/micro/go-micro/v2/transport"
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
		timeout:  u.opts.Timeout,
	}, nil
}

func (u *udpTransport) Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error) {

	var options transport.ListenOptions
	for _, o := range opts {
		o(&options)
	}
	var err error

	udpAddress, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenUDP("udp", udpAddress)

	if err != nil {
		return nil, err
	}
	return &udpListener{
		timeout:  u.opts.Timeout,
		sockets:  make(chan *udpSocket),
		exit:     make(chan bool),
		listener: l,
		pConn:    l,
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

//　　UDP : 1500 - IP(20) - UDP(8) = 1472(Bytes)
//Accept and handle a data package
func (u *udpListener) Accept(fn func(transport.Socket)) error {
	for {
		buf := make([]byte, defaultUDPMaxPackageLenth)
		//	rbuffer := ring.New(defaultUDPMaxPackageLenth),
		//conn, err := u.listener.Accept()
		//bytesLenth, fromAddr, err := u.listener.ReadFromUDP(buf)
		bytesLennth, fromAddr, err := u.pConn.ReadFrom(buf)

		if err != nil {
			u.errorChan <- err
		}
		select {
		case <-u.exit:
			return nil
		case <-u.errorChan:
			return nil
		default:
			sock := &udpSocket{
				timeout: u.opts.Timeout,
				conn:    u.listener,
				pConn:   u.listener,
				remote:  fromAddr.String(),
				local:   u.Addr(),
				encBuf: bufio.NewWriter(u.listener)
				exit:   make(chan bool)
			}
			go fn(sock)

		}
	}
}


func (u *udpListener) Addr() string {
	return u.listener.LocalAddr().String()
}

func (u *udpListener) Close() error {
	return u.listener.Close()
}
