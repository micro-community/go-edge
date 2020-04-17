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
	//p, err := net.ListenPacket("udp", addr)

	if err != nil {
		return nil, err
	}
	return &udpListener{
		timeout:  u.opts.Timeout,
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
	if de, ok := deFromContext(u.opts.Context); ok {
		u.dataExtractor = de
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

		buf := make([]byte, UDPServerRecvMaxLen)
		//	rbuffer := ring.New(defaultUDPMaxPackageLenth),
		n, fromAddr, err := u.pConn.ReadFrom(buf)
		// the n > 0 bytes returned before considering the error err.
		if n <= 0 {
			continue
		}
		if err != nil {
			u.errorChan <- struct{}{}
		}
		select {
		case <-u.exit:
			return nil
		case <-u.errorChan:
			return nil
		default:
			sock := &udpSocket{
				timeout: u.timeout,
				conn:    u.listener,
				pConn:   u.listener,
				dstAddr: fromAddr,
				remote:  fromAddr.String(),
				local:   u.Addr(),
				encBuf:  bufio.NewWriter(u.listener),
				exit:    make(chan bool),
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
