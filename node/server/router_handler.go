package server

import (
	"reflect"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/server"
)

type routingHandler struct {
	name    string
	handler interface{}
	opts    server.HandlerOptions
	typ     reflect.Type
}

//override the old
func newRoutingHandler(handler interface{}, opts ...server.HandlerOption) server.Handler {
	options := server.HandlerOptions{
		Metadata: make(map[string]map[string]string),
	}

	for _, o := range opts {
		o(&options)
	}

	typ := reflect.TypeOf(handler)
	hdlr := reflect.ValueOf(handler)
	name := reflect.Indirect(hdlr).Type().Name()

	return &routingHandler{
		name:    name,
		handler: handler,
		opts:    options,
		typ:     typ,
	}
}

func (r *routingHandler) Name() string {
	return r.name
}

func (r *routingHandler) Handler() interface{} {
	return r.handler
}

func (r *routingHandler) Options() server.HandlerOptions {
	return r.opts
}

func (r *routingHandler) Endpoints() []*registry.Endpoint {
	return nil
}
