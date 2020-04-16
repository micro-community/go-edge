// Package udp provides a udp transport
package udp

import (
	"bufio"
	"net"
	"sync"
	"time"

	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/transport"
)

//UDPServerRecvMaxLen Default UDP buffer len
const UDPServerRecvMaxLen = 1473

type udpTransport struct {
	opts      transport.Options
	listening chan struct{} // is closed when listen returns
}

type udpClient struct {
	dialOpts transport.DialOptions
	conn     net.Conn
	encBuf   *bufio.Writer
	timeout  time.Duration
}

type udpSocket struct {
	sync.RWMutex
	recv    chan *transport.Message
	send    chan *transport.Message
	conn    *net.UDPConn
	pConn   net.PacketConn
	encBuf  *bufio.Writer
	timeout time.Duration
	dstAddr *net.UDPAddr
	local   string
	remote  string
	exit    chan bool
}

type udpListener struct {
	sync.RWMutex
	timeout  time.Duration
	listener *net.UDPConn // current listener
	pConn    net.PacketConn
	//	sockets   chan *udpSocket
	errorChan chan struct{}
	exit      chan bool // sock exit
	closed    chan bool // listener exit
	opts      transport.ListenOptions
}

func init() {
	cmd.DefaultTransports["udp"] = NewTransport
}

//NewTransport Create a udp transport
func NewTransport(opts ...transport.Option) transport.Transport {
	var options transport.Options
	for _, o := range opts {
		o(&options)
	}
	return &udpTransport{opts: options}
}
