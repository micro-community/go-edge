package edge

import (
	"errors"
	"github.com/google/uuid"
	micro "github.com/micro/go-micro"
)

//Default Config
var (
	// For serving
	DefaultName    = "x-edge-node-srv"
	DefaultVersion = "latest"
	DefaultID      = uuid.New().String()
	DefaultAddress = ":8000"

	DefaultExtractor = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return -1, nil, errors.New("No Extractor Defined")
	}
)

// NewService returns a new web.Service
func NewService(opts ...Option) micro.Service {
	return newService(opts...)
}
