package edge

// Service is a edge service to connect to device/gw/controller/box
import (
	"bufio"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
)

// Service is a edge srv node running in edge process
type Service interface {
	Name() string
	//	service.Service
	Client() client.Client
	Server() server.Server
	Init(opts ...Option) error
	Options() Options
	//Handle(pattern string, handler http.Handler)
	//	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Run() error
	Stop() error
	String() string
}

//Option for edge
type Option func(o *Options)

//PackageExtractor for extract package from protocol
type PackageExtractor = bufio.SplitFunc

//service metadata
var (
	log = logger.NewHelper(logger.DefaultLogger).WithFields(map[string]interface{}{"service": "edge-node"})
)

// NewServer returns a new edge node server
func NewServer(opts ...Option) Service {
	return newService(opts...)
}
