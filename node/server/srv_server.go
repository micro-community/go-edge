package server

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	xmlc "github.com/micro-community/x-edge/node/codec"

	"github.com/micro/go-micro/codec"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"
	log "github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/util/socket"
)

type stubServer struct {
	router   *Routing
	opts     server.Options
	handlers map[string]server.Handler

	exit chan chan error
	sync.RWMutex
	// marks the serve as started
	started bool
	// graceful exit
	wg *sync.WaitGroup
}

// type stubRouter struct {
// 	h func(context.Context, Request, interface{}) error
// }

func newOption(opt ...server.Option) server.Options {
	opts := server.Options{
		Codecs:   make(map[string]codec.NewCodec),
		Metadata: map[string]string{},
	}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

//NewServer return a new custom rpc server
func NewServer(opts ...server.Option) server.Server {
	options := newOption(opts...)
	router := DefaultRouter()

	return &stubServer{
		opts:     options,
		router:   router,
		handlers: make(map[string]server.Handler),
		exit:     make(chan chan error),
		wg:       wait(options.Context),
	}
}

// ServeConn serves a single connection
func (s *stubServer) ServeConn(sock transport.Socket) {
	defer func() {
		// close socket
		sock.Close()

		if r := recover(); r != nil {
			log.Log("panic recovered: ", r)
			log.Log(string(debug.Stack()))
		}
	}()

	// multiplex the streams on a single socket by Micro-Stream
	var mtx sync.RWMutex
	sockets := make(map[string]*socket.Socket)

	for {
		var msg transport.Message
		if err := sock.Recv(&msg); err != nil {
			return
		}

		//as a key to  represent a session.
		id := sock.Local() + "-" + sock.Remote()

		// add to wait group if "wait" is opt-in
		if s.wg != nil {
			s.wg.Add(1)
		}

		// check we have an existing socket
		mtx.RLock()
		psock, ok := sockets[id]
		mtx.RUnlock()

		// got the socket
		if ok {
			// accept the message
			if err := psock.Accept(&msg); err != nil {
				// delete the socket
				mtx.Lock()
				delete(sockets, id)
				mtx.Unlock()
			}

			// done(1)
			if s.wg != nil {
				s.wg.Done()
			}

			// continue to the next message
			continue
		}

		// no socket was found
		psock = socket.New(id)
		psock.SetLocal(sock.Local())
		psock.SetRemote(sock.Remote())

		// load the socket
		psock.Accept(&msg)

		// save a new socket
		mtx.Lock()
		sockets[id] = psock
		mtx.Unlock()

		// process the outbound messages from the socket
		go func(id string, psock *socket.Socket) {
			defer psock.Close()

			for {
				// get the message from our internal handler/stream
				m := new(transport.Message)
				if err := psock.Process(m); err != nil {
					// delete the socket
					mtx.Lock()
					delete(sockets, id)
					mtx.Unlock()
					return
				}

				// send the message back over the socket
				if err := sock.Send(m); err != nil {
					return
				}
			}
		}(id, psock)

		msg.Header = map[string]string{}
		// set local/remote/codec for protocol
		msg.Header["Local"] = sock.Local()
		msg.Header["Remote"] = sock.Remote()
		msg.Header["Codec"] = xmlc.DefaultContentType

		msgCodec := s.newCodec(xmlc.DefaultContentType, psock)
		hdr := make(map[string]string)
		for k, v := range msg.Header {
			hdr[k] = v
		}

		// create new context with the metadata
		ctx := metadata.NewContext(context.Background(), hdr)

		// internal request
		rqst := &request{
			contentType: xmlc.DefaultContentType,
			codec:       msgCodec,
			header:      msg.Header,
			body:        msg.Body,
			socket:      psock,
		}

		// internal response
		resp := &response{
			header: make(map[string]string),
			socket: psock,
			codec:  msgCodec,
		}

		// serve the request in a go routine as this may be a stream
		go func(id string, psock *socket.Socket) {
			defer psock.Close()
			// serve the actual request using the request router
			if err := s.router.ServeRequest(ctx, rqst, resp); err != nil {
				log.Logf("unable to write error response: %v", err)
			}
			mtx.Lock()
			delete(sockets, id)
			mtx.Unlock()
			// signal we're done
			if s.wg != nil {
				s.wg.Done()
			}
		}(id, psock)

	}
}

//newCodec return codec for message.
func (s *stubServer) newCodec(contentType string, socket *socket.Socket) codec.Codec {
	if cf, ok := s.opts.Codecs[contentType]; ok {
		return newBuffCodec(socket, cf)
	}
	log.Errorf("Unsupported Content-Type: %s", contentType)
	//Default for xml
	return newBuffCodec(socket, xmlc.DefaultCodecs[contentType])
}

func (s *stubServer) Options() server.Options {
	return s.opts
}

func (s *stubServer) Init(opts ...server.Option) error {
	s.Lock()
	for _, opt := range opts {
		opt(&s.opts)
	}
	s.Unlock()
	return nil
}

func (s *stubServer) NewHandler(h interface{}, opts ...server.HandlerOption) server.Handler {
	return s.router.NewHandler(h, opts...)
}

func (s *stubServer) Handle(h server.Handler) error {
	s.Lock()
	defer s.Unlock()

	if err := s.router.Handle(h); err != nil {
		return err
	}

	s.handlers[h.Name()] = h

	return nil
}

func (s *stubServer) NewSubscriber(topic string, sb interface{}, opts ...server.SubscriberOption) server.Subscriber {
	return nil
}

func (s *stubServer) Subscribe(sb server.Subscriber) error {
	return nil
}

//Register useless here
func (s *stubServer) Register() error {
	return nil
}

//Deregister useless here
func (s *stubServer) Deregister() error {
	return nil
}

func (s *stubServer) Start() error {
	s.RLock()
	if s.started {
		s.RUnlock()
		return nil
	}
	s.RUnlock()

	config := s.opts

	// start listening on the transport
	ts, err := config.Transport.Listen(config.Address)
	if err != nil {
		return err
	}

	log.Logf("Transport [%s] Listening on %s", config.Transport.String(), ts.Addr())

	// swap address
	//s.Lock()
	//addr := s.rs.Options().Address
	//s.Options().Address = ts.Addr()
	//s.Unlock()

	exit := make(chan bool)

	go func() {
		for {
			// listen for connections
			err := ts.Accept(s.ServeConn)

			// TODO: listen for messages
			// msg := broker.Exchange(service).Consume()

			select {
			// check if we're supposed to exit
			case <-exit:
				return
			// check the error and backoff
			default:
				if err != nil {
					log.Logf("Accept error: %v", err)
					time.Sleep(time.Second)
					continue
				}
			}

			// no error just exit
			return
		}
	}()

	// mark the server as started
	s.Lock()
	s.started = true
	s.Unlock()

	return nil
}

func (s *stubServer) Stop() error {
	s.RLock()
	if !s.started {
		s.RUnlock()
		return nil
	}
	s.RUnlock()

	ch := make(chan error)
	//	s.exit <- ch

	var err error
	select {
	case err = <-ch:
		s.Lock()
		s.started = false
		s.Unlock()
	}

	return err
}

func (s *stubServer) String() string {
	return "stubServer"
}
