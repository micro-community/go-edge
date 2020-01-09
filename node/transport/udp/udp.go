// Package udp provides a udp transport
package udp

import (
	"bufio"
	"encoding/gob"
	"net"
	"time"

	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/transport"
)

type udpTransport struct {
	opts transport.Options
}

type udpClient struct {
	dialOpts transport.DialOptions
	conn     net.Conn
	enc      *gob.Encoder
	dec      *gob.Decoder
	encBuf   *bufio.Writer
	timeout  time.Duration
}

type udpSocket struct {
	conn *net.UDPConn
	// enc        *gob.Encoder
	// dec        *gob.Decoder
	encBuf     *bufio.Writer
	timeout    time.Duration
	packageBuf []byte
	packageLen int
	dstAddr    *net.UDPAddr
	closed     bool
}

//UDPServerRecvMaxLen Default UDP buffer len
const UDPServerRecvMaxLen = 2048

type udpListener struct {
	timeout  time.Duration
	listener *net.UDPConn
	opts     transport.ListenOptions
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

