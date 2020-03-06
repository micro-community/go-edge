package client

import (
	"context"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/codec"
	"github.com/micro/go-micro/v2/transport"
)

var (
	// DefaultPoolMaxStreams maximum streams on a connectioin
	// (20)
	DefaultPoolMaxStreams = 20

	// DefaultPoolMaxIdle maximum idle conns of a pool
	// (50)
	DefaultPoolMaxIdle = 50

	// DefaultMaxRecvMsgSize maximum message that client can receive
	// (1K).
	DefaultMaxRecvMsgSize = 1024

	// DefaultMaxSendMsgSize maximum message that client can send
	// (1K).
	DefaultMaxSendMsgSize = 1024
)

type codecsKey struct{}
type maxRecvMsgSizeKey struct{}
type maxSendMsgSizeKey struct{}
type nodeDialOptions struct{}
type nodeCallOptions struct{}

//Codec to be used to encode/decode requests for a given content type
func Codec(contentType string, c codec.NewCodec) client.Option {
	return func(o *client.Options) {
		o.Codecs[contentType] = c
	}
}

// MaxRecvMsgSize set the maximum size of message that client can receive.
func MaxRecvMsgSize(s int) client.Option {
	return func(o *client.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, maxRecvMsgSizeKey{}, s)
	}
}

// MaxSendMsgSize set the maximum size of message that client can send.
func MaxSendMsgSize(s int) client.Option {
	return func(o *client.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, maxSendMsgSizeKey{}, s)
	}
}

// CallOptions to be used to configure node call options
func CallOptions(opts ...client.CallOption) client.CallOption {
	return func(o *client.CallOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, nodeCallOptions{}, opts)
	}
}

// Transport to use for communication e.g http, rabbitmq, etc
func Transport(t transport.Transport) client.Option {
	return func(o *client.Options) {
		o.Transport = t
	}
}
