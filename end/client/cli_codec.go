package client

import (
	xmlc "github.com/micro-community/x-edge/end/codec"
	"github.com/micro-community/x-edge/end/iobuffer"
	"github.com/micro/go-micro/codec"
	raw "github.com/micro/go-micro/codec/bytes"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/transport/tcp"
)

var (
	defaultCodecs = map[string]codec.NewCodec{
		"application/xml": xmlc.NewCodec,
	}

	//DefaultContentType xml
	DefaultContentType = "application/xml"
)

type codecBuffer struct {
	socket transport.Socket
	codec  codec.Codec
	first  bool

	req *transport.Message //buffer the req msg
	buf *iobuffer.ReadWriteCloser
}

func newBuffCodec(client transport.Client, c codec.NewCodec, stream bool) codec.Codec {
	rwc := iobuffer.NewBuffer()

	cb := &codecBuffer{
		buf:   rwc,
		codec: c(rwc),
		//		req:    req,
		socket: client,
	}
	return cb
}

func (c *codecBuffer) ReadHeader(m *codec.Message, t codec.MessageType) error {

	var tm transport.Message

	// read message from transport
	if err := c.socket.Recv(&tm); err != nil {
		return errors.InternalServerError("go.micro.client.transport", err.Error())
	}

	c.buf.Reset()
	// write the body to the buffer
	if _, err := c.buf.Write(tm.Body); err != nil { //WriteWbuf
		//if _, err := c.buf.WriteWbuf(tm.Body); err != nil { //
		return err
	}

	// set headers from transport
	m.Header = tm.Header

	// read header
	err := c.codec.ReadHeader(m, t)

	// return header error
	if err != nil {
		return errors.InternalServerError("go.micro.client.codec", err.Error())
	}

	return nil
}

func (c *codecBuffer) ReadBody(b interface{}) error {
	// don't read empty body
	if len(c.req.Body) == 0 {
		return nil
	}
	// read raw data
	if v, ok := b.(*raw.Frame); ok {
		v.Data = c.buf.RBytes()
		return nil
	}

	if err := c.codec.ReadBody(b); err != nil {
		return errors.InternalServerError("client.codec", err.Error())
	}

	return nil
}

func (c *codecBuffer) Write(m *codec.Message, body interface{}) error {
	c.buf.Reset()

	// create header
	if m.Header == nil {
		m.Header = map[string]string{}
	}

	// copy original header
	for k, v := range c.req.Header {
		m.Header[k] = v
	}

	// if body is bytes Frame don't encode
	if body != nil {
		b, ok := body.(*raw.Frame)
		if ok {
			// set body
			m.Body = b.Data
			body = nil
		}
	}

	if len(m.Body) == 0 {
		// write to codec
		if err := c.codec.Write(m, body); err != nil {
			return errors.InternalServerError("go.micro.client.codec", err.Error())
		}
		// set body
		m.Body = c.buf.WBytes()
	}

	// create new transport message
	msg := transport.Message{
		Header: m.Header,
		Body:   m.Body,
	}
	// send the request
	if err := c.socket.Send(&msg); err != nil {
		return errors.InternalServerError("end.transport", err.Error())
	}
	return nil
}

func (c *codecBuffer) Close() error {
	c.buf.Close()
	c.codec.Close()
	return c.socket.Close()
}

func (c *codecBuffer) String() string {
	return "codecBuffer"
}
