// Package udp provides a udp transport
package udp

import (
	"bufio"
	"net"
	"sync"
	"time"

	nts "github.com/micro-community/x-edge/node/transport"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/transport"
)

func init() {
	cmd.DefaultTransports["udp"] = NewTransport
}

//UDPServerRecvMaxLen Default UDP buffer len
//　　UDP : 1500 - IP(20) - UDP(8) = 1472(Bytes)
const UDPServerRecvMaxLen = 1472

type udpTransport struct {
	opts          transport.Options
	dataExtractor nts.DataExtractor
	listening     chan struct{} // is closed when listen returns
}

type udpClient struct {
	dialOpts      transport.DialOptions
	conn          net.Conn
	pConn         net.PacketConn
	encBuf        *bufio.Writer
	timeout       time.Duration
	dataExtractor nts.DataExtractor
}

type udpSocket struct {
	sync.RWMutex
	recv          chan *transport.Message
	send          chan *transport.Message
	conn          *net.UDPConn
	pConn         net.PacketConn
	encBuf        *bufio.Writer
	timeout       time.Duration
	dstAddr       net.Addr
	local         string
	remote        string
	exit          chan bool
	dataExtractor nts.DataExtractor
	packageBuf    []byte
	packageLen    int
	closed        bool
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

//NewTransport Create a udp transport
func NewTransport(opts ...transport.Option) transport.Transport {
	var options transport.Options
	for _, o := range opts {
		o(&options)
	}
	return &udpTransport{opts: options}
}
