package server

import (
	"github.com/micro-community/x-edge/end/iobuffer"
	"github.com/micro/go-micro/codec"
	raw "github.com/micro/go-micro/codec/bytes"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-micro/util/socket"
)

type codecBuffer struct {
	socket transport.Socket
	codec  codec.Codec
	first  bool

	req *transport.Message //buffer the req msg
	buf *iobuffer.ReadWriteCloser
}

func newBuffCodec(sock *socket.Socket, c codec.NewCodec) codec.Codec {
	rwc := iobuffer.NewBuffer()

	cb := &codecBuffer{
		buf:   rwc,
		codec: c(rwc),
		//		req:    req,
		socket: sock,
	}
	return cb
}

func (c *codecBuffer) ReadHeader(r *codec.Message, t codec.MessageType) error {

	var tm transport.Message

	// read off the socket
	if err := c.socket.Recv(&tm); err != nil {
		return err
	}
	// reset the read buffer
	c.buf.Reset()

	// write the body to the buffer
	if _, err := c.buf.WriteRbuf(tm.Body); err != nil {
		return err
	}

	// the initial message
	m := codec.Message{
		Header: tm.Header,
		Body:   tm.Body,
	}

	// set req
	c.req = &tm

	// read header via codec
	if err := c.codec.ReadHeader(&m, codec.Request); err != nil {
		return err
	}

	// fallback for 0.14 and older
	if len(m.Endpoint) == 0 {
		m.Endpoint = m.Method
	}

	// set message
	*r = m

	return nil
}

func (c *codecBuffer) ReadBody(b interface{}) error {
	// don't read empty body
	if len(c.req.Body) == 0 {
		return nil
	}
	// read raw data
	if v, ok := b.(*raw.Frame); ok {
		v.Data = c.req.Body
		return nil
	}

	if pb, ok := b.(*codec.Message); ok {
		pb.Body = c.req.Body
		return nil
	}

	if err := c.codec.ReadBody(b); err != nil {
		return errors.InternalServerError("client.codec", err.Error())
	}

	return nil
}

func (c *codecBuffer) Write(r *codec.Message, b interface{}) error {
	c.buf.Reset()

	// create a new message
	m := &codec.Message{
		Target:   r.Target,
		Method:   r.Method,
		Endpoint: r.Endpoint,
		Id:       r.Id,
		Error:    r.Error,
		Type:     r.Type,
		Header:   r.Header,
	}

	if m.Header == nil {
		m.Header = map[string]string{}
	}

	// the body being sent
	var body []byte

	// is it a raw frame?
	if v, ok := b.(*raw.Frame); ok {
		body = v.Data
		// if we have encoded data just send it
	} else if len(r.Body) > 0 {
		body = r.Body
		// write the body to codec
	} else if err := c.codec.Write(m, b); err != nil {
		c.buf.Reset()
		// no body to write
		if err := c.codec.Write(m, nil); err != nil {
			return err
		}
	} else {
		// set the body
		body = c.buf.WBytes()
	}

	// Set content type if theres content
	if len(body) > 0 {
		m.Header["Content-Type"] = c.req.Header["Content-Type"]
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
