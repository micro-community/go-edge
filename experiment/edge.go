// Package experiment provides experiment edge based micro services
package experiment

import (
	"github.com/google/uuid"
	edge "github.com/micro-community/x-edge"
)

// Service is a web service with service discovery built in
type Service interface {
	Client() *edge.Client
	Init(opts ...Option) error
	Options() Options
	Handle(pattern string, handler edge.Handler)
	HandleFunc(pattern string, handler func(edge.ResponseWriter, *edge.Request))
	Run() error
}

type Option func(o *Options)

//Default Config
var (
	// For serving
	DefaultName    = "x-edge"
	DefaultVersion = "latest"
	DefaultID      = uuid.New().String()
	DefaultAddress = ":0"
)

// NewService returns a new web.Service
func NewService(opts ...Option) Service {
	return newService(opts...)
}
