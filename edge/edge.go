package edge

import (
	"bufio"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/server"
)

// Service is a edge srv node running in edge process
type Service interface {
	Name() string
	Client() client.Client
	Server() server.Server
	Init(opts ...Option) error
	Options() Options
	Start() error
	Stop() error
	String() string
}

//Option for edge
type Option func(o *Options)

//PackageExtractor for extract package from protocol
type PackageExtractor = bufio.SplitFunc

// NewServer returns a new edge node server
func NewServer(opts ...Option) Service {
	return newService(opts...)
}
