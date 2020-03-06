package router

import (
	context "context"

	"github.com/micro/go-micro/v2/codec"
	server "github.com/micro/go-micro/v2/server"
)

//some constant variable
var (
	DefaultServerName = "ProtocolService"
)

//ProtocolHandler for Protocol Handling
type ProtocolHandler interface {
	Event(ctx context.Context, in *codec.Message, out *codec.Message) error
}

//RegisterProtocolHandler to router
func RegisterProtocolHandler(s server.Server, hdlr ProtocolHandler, opts ...server.HandlerOption) error {
	return s.Handle(s.NewHandler(hdlr, opts...))
}
