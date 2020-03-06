package client

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/micro/go-micro/v2/transport"
	"github.com/micro/go-micro/v2/util/log"

	xmlc "github.com/micro-community/x-edge/node/codec"
	"github.com/micro-community/x-edge/node/stream"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/pool"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/codec"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
)

type nodeClient struct {
	once sync.Once
	opts client.Options
	pool pool.Pool
	seq  uint64
}

func newOption(options ...client.Option) client.Options {

	opts := client.Options{
		Codecs: make(map[string]codec.NewCodec),
		CallOptions: client.CallOptions{
			Backoff:        client.DefaultBackoff,
			Retry:          client.DefaultRetry,
			Retries:        client.DefaultRetries,
			RequestTimeout: client.DefaultRequestTimeout,
			DialTimeout:    transport.DefaultDialTimeout,
		},
		PoolSize: client.DefaultPoolSize,
		PoolTTL:  client.DefaultPoolTTL,
	}

	for _, o := range options {
		o(&opts)
	}
	return opts
}

//NewClient return a new custom rpc client
func NewClient(opts ...client.Option) client.Client {

	options := newOption(opts...)

	p := pool.NewPool(
		pool.Size(options.PoolSize),
		pool.TTL(options.PoolTTL),
		pool.Transport(options.Transport),
	)

	sc := &nodeClient{
		once: sync.Once{},
		opts: options,
		pool: p,
		seq:  0,
	}

	c := client.Client(sc)

	// wrap in reverse
	for i := len(options.Wrappers); i > 0; i-- {
		c = options.Wrappers[i-1](c)
	}

	return c
}

func (c *nodeClient) newCodec(contentType string, client transport.Client, stream bool) codec.Codec {
	if cf, ok := c.opts.Codecs[contentType]; ok {
		return newBuffCodec(client, cf, stream)
	}
	log.Infof("Unsupported Content-Type: %s", contentType)
	return newBuffCodec(client, xmlc.DefaultCodecs[contentType], stream)

}

func (c *nodeClient) call(ctx context.Context, req client.Request, resp interface{}, opts client.CallOptions) error {

	address := ctx.Value("remote").(string)

	msg := &transport.Message{
		Header: make(map[string]string),
	}

	md, ok := metadata.FromContext(ctx)
	if ok {
		for k, v := range md {
			msg.Header[k] = v
		}
	}

	// set timeout in nanoseconds
	msg.Header["Timeout"] = fmt.Sprintf("%d", opts.RequestTimeout)
	// set the content type for the request
	msg.Header["Content-Type"] = req.ContentType()
	// set the accept header
	msg.Header["Accept"] = req.ContentType()

	con, err := c.pool.Get(address, transport.WithTimeout(opts.DialTimeout))
	if err != nil {
		return errors.InternalServerError("go.micro.raw.client", "connection error: %v", err)
	}

	seq := atomic.LoadUint64(&c.seq)
	atomic.AddUint64(&c.seq, 1)

	//if this is a file ,it should be a stream,now we just ignore it.
	msgCodec := c.newCodec(xmlc.DefaultContentType, con, false)

	rsp := &response{
		socket: con,
		codec:  msgCodec,
	}

	releaseOP := func(err error) { c.pool.Release(con, err) }

	stream := stream.NewClientStrem(ctx, seq, req, rsp, msgCodec, releaseOP)

	defer stream.Close()

	//wait for error response
	ch := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				ch <- errors.InternalServerError("go.micro.client", "panic recovered: %v", r)
			}
		}()

		// send request
		if err := stream.Send(req.Body()); err != nil {
			ch <- err
			return
		}

		// recv request
		if err := stream.Recv(resp); err != nil {
			ch <- err
			return
		}

		// success
		ch <- nil
	}()

	var grr error

	select {
	case err := <-ch:
		grr = err
		return err
	case <-ctx.Done():
		grr = errors.Timeout("go.micro.client", fmt.Sprintf("%v", ctx.Err()))
	}

	return grr
}

func (c *nodeClient) stream(ctx context.Context, req client.Request, opts client.CallOptions) (client.Stream, error) {

	//address := ctx.Value("target-service").(string)
	//address := node.Address

	address := ctx.Value("remote").(string)

	msg := &transport.Message{
		Header: make(map[string]string),
	}

	md, ok := metadata.FromContext(ctx)
	if ok {
		for k, v := range md {
			msg.Header[k] = v
		}
	}

	// set timeout in nanoseconds
	msg.Header["Timeout"] = fmt.Sprintf("%d", opts.RequestTimeout)
	// set the content type for the request
	msg.Header["Content-Type"] = req.ContentType()
	// set the accept header
	msg.Header["Accept"] = req.ContentType()

	// set old codecs

	dOpts := []transport.DialOption{
		transport.WithStream(),
	}

	if opts.DialTimeout >= 0 {
		dOpts = append(dOpts, transport.WithTimeout(opts.DialTimeout))
	}

	con, err := c.pool.Get(address, dOpts...)
	if err != nil {
		return nil, errors.InternalServerError("go.micro.client", "connection error: %v", err)
	}

	// increment the sequence number
	seq := atomic.LoadUint64(&c.seq)
	atomic.AddUint64(&c.seq, 1)

	//if this is a file ,it should be a stream,now we just ignore it.
	msgCodec := c.newCodec(xmlc.DefaultContentType, con, false)

	rsp := &response{
		socket: con,
		codec:  msgCodec,
	}

	releaseOP := func(err error) { c.pool.Release(con, err) }

	stream := stream.NewClientStrem(ctx, seq, req, rsp, msgCodec, releaseOP)

	// wait for error response
	ch := make(chan error, 1)

	go func() {
		// send the first message
		ch <- stream.Send(req.Body())
	}()

	var grr error

	select {
	case err := <-ch:
		grr = err
	case <-ctx.Done():
		grr = errors.Timeout("go.micro.client", fmt.Sprintf("%v", ctx.Err()))
	}

	if grr != nil {
		// set the error

		stream.Lock()
		stream.SetError(grr)
		stream.Unlock()

		// close the stream
		stream.Close()
		return nil, grr
	}

	return stream, nil
}

func (c *nodeClient) Init(opts ...client.Option) error {
	size := c.opts.PoolSize
	ttl := c.opts.PoolTTL
	tr := c.opts.Transport

	for _, o := range opts {
		o(&c.opts)
	}

	// update pool configuration if the options changed
	if size != c.opts.PoolSize || ttl != c.opts.PoolTTL || tr != c.opts.Transport {
		// close existing pool
		c.pool.Close()
		// create new pool
		c.pool = pool.NewPool(
			pool.Size(c.opts.PoolSize),
			pool.TTL(c.opts.PoolTTL),
			pool.Transport(c.opts.Transport),
		)
	}

	return nil
}

func (c *nodeClient) Options() client.Options {
	return c.opts
}

//we will use static selector
func (c *nodeClient) next(request client.Request, opts client.CallOptions) (selector.Next, error) {

	service := request.Service()

	// get next nodes from the selector
	next, err := c.opts.Selector.Select(service, opts.SelectOptions...)
	if err != nil {
		if err == selector.ErrNotFound {
			return nil, errors.InternalServerError("go.micro.client", "service %s: %s", service, err.Error())
		}
		return nil, errors.InternalServerError("go.micro.client", "error selecting %s node: %s", service, err.Error())
	}

	return next, nil
}

func (c *nodeClient) Call(ctx context.Context, request client.Request, response interface{}, opts ...client.CallOption) error {
	// make a copy of call opts
	callOpts := c.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	next, err := c.next(request, callOpts)
	if err != nil {
		return err
	}

	var concelFunc context.CancelFunc
	//check if we already have a deadline
	d, ok := ctx.Deadline()
	if !ok {
		// no deadline so we create a new one
		ctx, concelFunc = context.WithTimeout(ctx, callOpts.RequestTimeout)
		defer concelFunc()
	} else {
		// got a deadline so no need to setup context
		// but we need to set the timeout we pass along
		opt := client.WithRequestTimeout(d.Sub(time.Now()))
		opt(&callOpts)
	}

	// should we noop right here?
	select {
	case <-ctx.Done():
		return errors.Timeout("go.micro.client", fmt.Sprintf("%v", ctx.Err()))
	default:
	}

	// make copy of call method
	rcall := c.call

	// return errors.New("go.micro.client", "request timeout", 408)
	call := func(i int) error {
		// call backoff first. Someone may want an initial start delay
		t, err := callOpts.Backoff(ctx, request, i)
		if err != nil {
			return errors.InternalServerError("go.micro.client", "backoff error: %v", err.Error())
		}

		// only sleep if greater than 0
		if t.Seconds() > 0 {
			time.Sleep(t)
		}

		// select target service
		node, err := next()
		service := request.Service()
		if err != nil {
			if err == selector.ErrNotFound {
				return errors.InternalServerError("go.micro.client", "service %s: %s", service, err.Error())
			}
			return errors.InternalServerError("go.micro.client", "error getting next %s node: %s", service, err.Error())
		}

		// make the call
		err = rcall(ctx, request, response, callOpts)
		c.opts.Selector.Mark(service, node, err)

		return err
	}

	ch := make(chan error, callOpts.Retries+1)
	var gerr error

	for i := 0; i <= callOpts.Retries; i++ {
		go func(i int) {
			ch <- call(i)
		}(i)

		select {
		case <-ctx.Done():
			return errors.Timeout("go.micro.client", fmt.Sprintf("call timeout: %v", ctx.Err()))
		case err := <-ch:
			// if the call succeeded lets bail early
			if err == nil {
				return nil
			}

			retry, rerr := callOpts.Retry(ctx, request, i, err)
			if rerr != nil {
				return rerr
			}

			if !retry {
				return err
			}

			gerr = err
		}
	}

	return gerr
}

func (c *nodeClient) Stream(ctx context.Context, request client.Request, opts ...client.CallOption) (client.Stream, error) {
	// make a copy of call opts
	callOpts := c.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	next, err := c.next(request, callOpts)
	if err != nil {
		return nil, err
	}

	// should we noop right here?
	select {
	case <-ctx.Done():
		return nil, errors.Timeout("go.micro.client", fmt.Sprintf("%v", ctx.Err()))
	default:
	}

	call := func(i int) (client.Stream, error) {
		// call backoff first. Someone may want an initial start delay
		t, err := callOpts.Backoff(ctx, request, i)
		if err != nil {
			return nil, errors.InternalServerError("go.micro.client", "backoff error: %v", err.Error())
		}

		// only sleep if greater than 0
		if t.Seconds() > 0 {
			time.Sleep(t)
		}

		node, err := next()
		service := request.Service()
		if err != nil {
			if err == selector.ErrNotFound {
				return nil, errors.InternalServerError("go.micro.client", "service %s: %s", service, err.Error())
			}
			return nil, errors.InternalServerError("go.micro.client", "error getting next %s node: %s", service, err.Error())
		}

		stream, err := c.stream(ctx, request, callOpts)
		c.opts.Selector.Mark(service, node, err)
		return stream, err
	}

	type response struct {
		stream client.Stream
		err    error
	}

	ch := make(chan response, callOpts.Retries+1)
	var grr error

	for i := 0; i <= callOpts.Retries; i++ {
		go func(i int) {
			s, err := call(i)
			ch <- response{s, err}
		}(i)

		select {
		case <-ctx.Done():
			return nil, errors.Timeout("go.micro.client", fmt.Sprintf("call timeout: %v", ctx.Err()))
		case rsp := <-ch:
			// if the call succeeded lets bail early
			if rsp.err == nil {
				return rsp.stream, nil
			}

			retry, rerr := callOpts.Retry(ctx, request, i, rsp.err)
			if rerr != nil {
				return nil, rerr
			}

			if !retry {
				return nil, rsp.err
			}

			grr = rsp.err
		}
	}

	return nil, grr
}

func (c *nodeClient) Publish(ctx context.Context, msg client.Message, opts ...client.PublishOption) error {
	return nil
}

func (c *nodeClient) NewMessage(topic string, message interface{}, opts ...client.MessageOption) client.Message {
	return newMessage(topic, message, c.opts.ContentType, opts...)
}

func (c *nodeClient) NewRequest(service, method string, request interface{}, reqOpts ...client.RequestOption) client.Request {
	return newRequest(service, method, request, c.opts.ContentType, reqOpts...)
}

func (c *nodeClient) String() string {
	return "nodeclient"
}
