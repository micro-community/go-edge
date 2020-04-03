package edge

// Service is a edge service to connect to device/gw/controller/box
import (
	"errors"
	//here should edge internal logic
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
)

// Service is a web service with service discovery built in
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
	ErrNoExtractorDefined = errors.New("No Extractor Defined")

	DefaultExtractor = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return -1, nil, ErrNoExtractorDefined
	}

	log = logger.NewHelper(logger.DefaultLogger).WithFields(map[string]interface{}{"service": "edge"})
)

// NewService returns a new web.Service
func NewService(opts ...Option) Service {
	return newService(opts...)
}
