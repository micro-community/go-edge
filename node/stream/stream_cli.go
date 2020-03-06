package stream

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/codec"
)

const (
	lastStreamResponseError = "EOS"
)

// Implements the Streamer interface
type streamClient struct {
	sync.RWMutex
	id       string
	closed   chan bool
	err      error
	request  client.Request
	response client.Response
	codec    codec.Codec
	context  context.Context

	// signal whether we should send EOS
	sendEOS bool

	// release releases the connection back to the pool
	release func(err error)
}

//SyncClientStream return a sync stream client
type SyncClientStream interface {
	client.Stream
	Lock()
	Unlock()
	SetError(err error)
}

//NewClientStrem return client stream
func NewClientStrem(ctx context.Context, seq uint64, req client.Request, rsp client.Response, msgCodec codec.Codec, r func(err error)) SyncClientStream {
	s := &streamClient{
		id:       fmt.Sprintf("%v", seq),
		request:  req,
		response: rsp,
		codec:    msgCodec,
		closed:   make(chan bool),
		sendEOS:  false,
		release:  r,
	}

	return s
}

func (r *streamClient) isClosed() bool {
	select {
	case <-r.closed:
		return true
	default:
		return false
	}
}

func (r *streamClient) Context() context.Context {
	return r.context
}

func (r *streamClient) Request() client.Request {
	return r.request
}

func (r *streamClient) Response() client.Response {
	return r.response
}

func (r *streamClient) Send(msg interface{}) error {
	r.Lock()
	defer r.Unlock()

	if r.isClosed() {
		r.err = errShutdown
		return errShutdown
	}

	req := codec.Message{
		Id:       r.id,
		Target:   r.request.Service(),
		Method:   r.request.Method(),
		Endpoint: r.request.Endpoint(),
		Type:     codec.Request,
	}

	if err := r.codec.Write(&req, msg); err != nil {
		r.err = err
		return err
	}

	return nil
}

func (r *streamClient) Recv(msg interface{}) error {
	r.Lock()
	defer r.Unlock()

	if r.isClosed() {
		r.err = errShutdown
		return errShutdown
	}

	var resp codec.Message

	if err := r.codec.ReadHeader(&resp, codec.Response); err != nil {
		if err == io.EOF && !r.isClosed() {
			r.err = io.ErrUnexpectedEOF
			return io.ErrUnexpectedEOF
		}
		r.err = err
		return err
	}

	switch {
	case len(resp.Error) > 0:
		// We've got an error response. Give this to the request;
		// any subsequent requests will get the ReadResponseBody
		// error if there is one.
		if resp.Error != lastStreamResponseError {
			r.err = errors.New(resp.Error)
		} else {
			r.err = io.EOF
		}
		if err := r.codec.ReadBody(nil); err != nil {
			r.err = err
		}
	default:
		if err := r.codec.ReadBody(msg); err != nil {
			r.err = err
		}
	}

	return r.err
}

func (r *streamClient) Error() error {
	r.RLock()
	defer r.RUnlock()
	return r.err
}

func (r *streamClient) Close() error {
	select {
	case <-r.closed:
		return nil
	default:
		close(r.closed)

		// send the node of stream message
		if r.sendEOS {
			// no need to check for error
			r.codec.Write(&codec.Message{
				Id:       r.id,
				Target:   r.request.Service(),
				Method:   r.request.Method(),
				Endpoint: r.request.Endpoint(),
				Type:     codec.Error,
				Error:    lastStreamResponseError,
			}, nil)
		}

		err := r.codec.Close()

		// release the connection
		r.release(r.Error())

		// return the codec error
		return err
	}
}

func (r *streamClient) SetError(err error) {
	r.err = err
}
