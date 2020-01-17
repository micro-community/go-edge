package client

import (
	"github.com/micro/go-micro/codec"
	"github.com/micro/go-micro/transport"
)

type response struct {
	header map[string]string
	body   []byte
	socket transport.Socket
	codec  codec.Codec
}

func (r *response) Codec() codec.Reader {
	return r.codec
}

func (r *response) Header() map[string]string {
	return r.header
}

func (r *response) Read() ([]byte, error) {
	var msg transport.Message

	if err := r.socket.Recv(&msg); err != nil {
		return nil, err
	}

	// set internals
	r.header = msg.Header
	r.body = msg.Body

	return msg.Body, nil
}
