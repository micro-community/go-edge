package client

import (
	"github.com/micro/go-micro/client"

	"github.com/micro/go-micro/codec"
)

type request struct {
	service     string
	method      string
	endpoint    string
	contentType string
	codec       codec.Codec
	body        interface{}
	opts        client.RequestOptions
}

func newRequest(service, endpoint string, req interface{}, contentType string, reqOpts ...client.RequestOption) client.Request {
	var opts client.RequestOptions

	for _, o := range reqOpts {
		o(&opts)
	}

	// set the content-type specified
	if len(opts.ContentType) > 0 {
		contentType = opts.ContentType
	}

	return &request{
		service:     service,
		method:      endpoint,
		endpoint:    endpoint,
		body:        req,
		contentType: contentType,
		opts:        opts,
	}
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

func (r *request) Body() interface{} {
	return r.body
}

func (r *request) Codec() codec.Writer {
	return r.codec
}

func (r *request) Stream() bool {
	return r.opts.Stream
}
