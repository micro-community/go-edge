package server

import (
	"bytes"
	"context"

	"github.com/micro/go-micro/v2/codec"
	"github.com/micro/go-micro/v2/transport"
	"github.com/micro/go-micro/v2/util/buf"
)

type request struct {
	service     string
	method      string
	endpoint    string
	contentType string
	socket      transport.Socket
	codec       codec.Codec
	header      map[string]string
	body        []byte
	rawBody     interface{}
	stream      bool
	first       bool
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

func (r *request) Codec() codec.Reader {
	return r.codec
}

func (r *request) ContentType() string {
	return r.contentType
}

func (r *request) Service() string {
	return r.service
}

func (r *request) Method() string {
	return r.method
}

func (r *request) Endpoint() string {
	return r.endpoint
}

func (r *request) Header() map[string]string {
	return r.header
}

func (r *request) Body() interface{} {
	return r.rawBody
}

func (r *request) Read() ([]byte, error) {
	// got a body
	if r.first {
		b := r.body
		r.first = false
		return b, nil
	}

	var msg transport.Message
	err := r.socket.Recv(&msg)
	if err != nil {
		return nil, err
	}
	r.header = msg.Header

	return msg.Body, nil
}

func (r *request) Stream() bool {
	return r.stream
}

type rpcMessage struct {
	topic       string
	contentType string
	payload     interface{}
	header      map[string]string
	body        []byte
	codec       codec.NewCodec
}

func (r *rpcMessage) ContentType() string {
	return r.contentType
}

func (r *rpcMessage) Topic() string {
	return r.topic
}

func (r *rpcMessage) Payload() interface{} {
	return r.payload
}

func (r *rpcMessage) Header() map[string]string {
	return r.header
}

func (r *rpcMessage) Body() []byte {
	return r.body
}

func (r *rpcMessage) Codec() codec.Reader {
	b := buf.New(bytes.NewBuffer(r.body))
	return r.codec(b)
}
