package service

// Service is a edge service to connect to device/gw/controller/box
import (
	"errors"

	"net/http" //here should edge internal logic

	"github.com/micro/go-micro/v2/logger"
)

// Service is a web service with service discovery built in
type Service interface {
	Client() *http.Client
	Init(opts ...Option) error
	Options() Options
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Run() error
}

//Option for edge
type Option func(o *Options)

//service metadata
var (
	DefaultName    = "x-edge-srv"
	DefaultAddress = ":8000"

	DefaultExtractor = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return -1, nil, errors.New("No Extractor Defined")
	}

	log = logger.NewHelper(logger.DefaultLogger).WithFields(map[string]interface{}{"service": "edge"})
)

// NewService returns a new web.Service
func NewService(opts ...Option) Service {
	return newService(opts...)
}
