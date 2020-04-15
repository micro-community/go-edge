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

type udpTransport struct {
	opts transport.Options
}

type udpClient struct {
	dialOpts transport.DialOptions
	conn     net.Conn
	encBuf   *bufio.Writer
	timeout  time.Duration
}

type udpSocket struct {
	sync.RWMutex
	recv chan *transport.Message
	send chan *transport.Message
	conn *net.UDPConn
	//	encBuf  *bufio.Writer
	timeout time.Duration
	//	packageBuf []byte
	//	packageLen int
	dstAddr *net.UDPAddr
	local   string
	remote  string
	exit    chan bool
}

//UDPServerRecvMaxLen Default UDP buffer len
const UDPServerRecvMaxLen = 2048

type udpListener struct {
	sync.RWMutex
	timeout  time.Duration
	listener *net.UDPConn
	conn     chan *udpSocket
	// sock exit
	exit chan bool
	// listener exit
	listenerexit chan bool
	opts         transport.ListenOptions
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
