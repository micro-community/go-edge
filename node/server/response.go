package server

import (
	"github.com/micro/go-micro/v2/codec"
	"github.com/micro/go-micro/v2/transport"
)

type response struct {
	header map[string]string
	socket transport.Socket
	codec  codec.Codec
}

func (r *response) Codec() codec.Writer {
	return r.codec
}

func (r *response) WriteHeader(hdr map[string]string) {
	for k, v := range hdr {
		r.header[k] = v
	}
}

func (r *response) Write(b []byte) error {

	return r.socket.Send(&transport.Message{
		Header: r.header,
		Body:   b,
	})
}
