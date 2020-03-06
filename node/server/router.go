package server

import (
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/micro/go-micro/v2/codec"
	"github.com/micro/go-micro/v2/server"
)

//Routing Config
type Routing struct {
	name       string
	mu         sync.Mutex // protects the serviceMap
	serviceMap map[string]*service
	reqLock    sync.Mutex       // protects freeReq
	respLock   sync.Mutex       // protects freeResp
	freeReq    *routingRequest  //request linklist
	freeResp   *routingResponse //response linklist
}

type routingRequest struct {
	msg  *codec.Message
	next *routingRequest // for free list in Server
}

type routingResponse struct {
	msg  *codec.Message
	next *routingResponse // for free list in Server
}

//DefaultRouter Return a new Router
func DefaultRouter() *Routing {
	router := &Routing{
		serviceMap: make(map[string]*service),
		mu:         sync.Mutex{},
	}

	return router
}

func (router *Routing) getRequest() *routingRequest {
	router.reqLock.Lock()
	req := router.freeReq
	if req == nil {
		req = new(routingRequest)
	} else {
		router.freeReq = req.next
		*req = routingRequest{}
	}
	router.reqLock.Unlock()
	return req
}

func (router *Routing) freeRequest(req *routingRequest) {
	router.reqLock.Lock()
	req.next = router.freeReq
	router.freeReq = req
	router.reqLock.Unlock()
}

func (router *Routing) getResponse() *routingResponse {
	router.respLock.Lock()
	resp := router.freeResp
	if resp == nil {
		resp = new(routingResponse)
	} else {
		router.freeResp = resp.next
		*resp = routingResponse{}
	}
	router.respLock.Unlock()
	return resp
}

func (router *Routing) freeResponse(resp *routingResponse) {
	router.respLock.Lock()
	resp.next = router.freeResp
	router.freeResp = resp
	router.respLock.Unlock()
}

func (router *Routing) readRequest(rqst server.Request) (service *service, mtype *methodType, req *routingRequest, argv, replyv reflect.Value, keepReading bool, err error) {

	codecBuffer := rqst.Codec() //codecBuffer codec
	req = router.getRequest()
	msg := &codec.Message{}
	msg.Type = codec.Request
	req.msg = msg

	err = codecBuffer.ReadHeader(msg, msg.Type)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
		}
		err = errors.New("router cannot decode request: " + err.Error())

		return
	}
	serviceName := strings.ToUpper(msg.Target)
	// Look up the request.
	router.mu.Lock()
	service = router.serviceMap[serviceName]
	router.mu.Unlock()

	if service == nil {
		err = errors.New(" can't find service: " + serviceName)
		return
	}
	methodName := strings.ToUpper(msg.Method)
	// find the matched method
	for methodNameInMap, targetMethod := range service.method {
		if strings.HasPrefix(methodNameInMap, methodName) {
			mtype = targetMethod
			break
		}
	}

	if mtype == nil {
		err = errors.New("can't find target method " + methodName)
	}

	// is it a streaming request? then we don't read the body
	if mtype.stream {
		codecBuffer.ReadBody(nil)
		return
	}

	// Decode the argument value.
	argIsValue := false // if true, need to indirect before calling.
	if mtype.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(mtype.ArgType.Elem())
	} else {
		argv = reflect.New(mtype.ArgType)
		argIsValue = true
	}
	// argv guaranteed to be a pointer now.
	if err = codecBuffer.ReadBody(argv.Interface()); err != nil {
		return
	}
	if argIsValue {
		argv = argv.Elem()
	}

	if !mtype.stream {
		replyv = reflect.New(mtype.ReplyType.Elem())
	}

	err = nil
	return
}

func (r *Routing) ProcessMessage(ctx context.Context, msg server.Message) error {
	return nil
}

//ServeRequest Serve requesting from controller
func (router *Routing) ServeRequest(ctx context.Context, rqst server.Request, rsp server.Response) error {
	sending := new(sync.Mutex)
	service, mtype, req, argv, replyv, keepReading, err := router.readRequest(rqst)
	//Here will receiving all request messages.
	if err != nil {
		if !keepReading {
			return err
		}
		// send a response if we actually managed to read a header.
		if req != nil {
			router.freeRequest(req)
		}
		return err
	}
	return service.call(ctx, router, sending, mtype, req, argv, replyv, rsp.Codec())
}

func (router *Routing) sendResponse(sending sync.Locker, req *routingRequest, reply interface{}, cc codec.Writer, last bool) error {
	msg := new(codec.Message)
	msg.Type = codec.Response
	resp := router.getResponse()
	resp.msg = msg

	resp.msg.Id = req.msg.Id
	sending.Lock()
	err := cc.Write(resp.msg, reply)
	sending.Unlock()
	router.freeResponse(resp)
	return err
}

//Handle function register
func (router *Routing) Handle(h server.Handler) error {
	router.mu.Lock()
	defer router.mu.Unlock()
	if router.serviceMap == nil {
		router.serviceMap = make(map[string]*service) //name and service map
	}

	if len(h.Name()) == 0 {
		return errors.New("rpc.Handle: handler has no name")
	}

	rcvr := h.Handler()
	s := new(service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)

	// check name
	if _, present := router.serviceMap[h.Name()]; present {
		return errors.New("router Handle: service already defined: " + h.Name())
	}

	s.name = strings.ToUpper(h.Name()) // this is "FunHandler"
	s.method = make(map[string]*methodType)

	// Install the methods
	for m := 0; m < s.typ.NumMethod(); m++ {
		method := s.typ.Method(m)
		if mt := prepareMethod(method); mt != nil {
			s.method[method.Name] = mt
		}
	}

	// Check there are methods
	if len(s.method) == 0 {
		return errors.New("rpc Register: type " + s.name + " has no exported methods of suitable type")
	}

	// save handler
	router.serviceMap[s.name] = s
	return nil
}

//NewHandler Add new Handler for target service
func (router *Routing) NewHandler(h interface{}, opts ...server.HandlerOption) server.Handler {
	return newRoutingHandler(h, opts...)
}
