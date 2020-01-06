package stream

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/micro/go-micro/codec"
	"github.com/micro/go-micro/server"
)

var (
	errorLastStreamResponse = errors.New("EOS")
	errShutdown             = errors.New("connection is shut down")
)

// Implements the Streamer interface
type streamServer struct {
	sync.RWMutex
	id       string
	closed   bool
	err      error
	request  server.Request
	response server.Response
	codec    codec.Codec
	context  context.Context
}

//NewServerStrem return a server stream object
func NewServerStrem(ctx context.Context, seq uint) server.Stream {

	s := &streamServer{}

	return s
}

func (r *streamServer) Context() context.Context {
	return r.context
}

func (r *streamServer) Request() server.Request {
	return r.request
}

func (r *streamServer) Send(msg interface{}) error {
	r.Lock()
	defer r.Unlock()

	resp := codec.Message{
		Target:   r.request.Service(),
		Method:   r.request.Method(),
		Endpoint: r.request.Endpoint(),
		ID:       r.id,
		Type:     codec.Response,
	}

	if err := r.codec.Write(&resp, msg); err != nil {
		r.err = err
	}

	return nil
}

func (r *streamServer) Recv(msg interface{}) error {
	r.Lock()
	defer r.Unlock()

	req := new(codec.Message)
	req.Type = codec.Request

	if err := r.codec.ReadHeader(req, req.Type); err != nil {
		// discard body
		r.codec.ReadBody(nil)
		r.err = err
		return err
	}

	// check the error
	if len(req.Error) > 0 {
		// Check the client closed the stream
		switch req.Error {
		case errorLastStreamResponse.Error():
			// discard body
			r.codec.ReadBody(nil)
			r.err = io.EOF
			return io.EOF
		default:
			return errors.New(req.Error)
		}
	}

	// we need to stay up to date with sequence numbers
	r.id = req.ID
	if err := r.codec.ReadBody(msg); err != nil {
		r.err = err
		return err
	}

	return nil
}

func (r *streamServer) Error() error {
	r.RLock()
	defer r.RUnlock()
	return r.err
}

func (r *streamServer) Close() error {
	r.Lock()
	defer r.Unlock()
	r.closed = true
	return r.codec.Close()
}
