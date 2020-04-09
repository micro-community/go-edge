package edge

// Service is a edge service to connect to device/gw/controller/box
import (
	"errors"
	//here should edge internal logic
	"github.com/google/uuid"
	nserver "github.com/micro-community/x-edge/node/server"
	"github.com/micro-community/x-edge/node/transport/tcp"
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

//service metadata
var (
	DefaultName           = "x-edge-node"
	DefaultVersion        = "latest"
	DefaultAddress        = ":8000"
	DefaultID             = uuid.New().String()
	DefaultServer         = nserver.NewServer()
	DefaultTransport      = tcp.NewTransport()
	ErrNoExtractorDefined = errors.New("No Extractor Defined")

	DefaultExtractor = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return -1, nil, ErrNoExtractorDefined
	}

	log = logger.NewHelper(logger.DefaultLogger).WithFields(map[string]interface{}{"service": "edge"})
)

// NewServer returns a new edge node server
func NewServer(opts ...Option) Service {
	return newService(opts...)
}
