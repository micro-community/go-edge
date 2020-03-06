package server

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/micro/go-micro/v2/codec"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/util/log"
)

var (

	// Precompute the reflect type for error. Can't use error directly
	// because Typeof takes an empty interface value. This is annoying.
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
)

type methodType struct {
	sync.Mutex  // protects counters
	method      reflect.Method
	ArgType     reflect.Type
	ReplyType   reflect.Type
	ContextType reflect.Type
	stream      bool
}

type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

func (m *methodType) prepareContext(ctx context.Context) reflect.Value {
	if contextv := reflect.ValueOf(ctx); contextv.IsValid() {
		return contextv
	}
	return reflect.Zero(m.ContextType)
}

func (s *service) call(ctx context.Context, router *Routing, sending *sync.Mutex, mtype *methodType, req *routingRequest, argv, replyv reflect.Value, cc codec.Writer) error {
	defer router.freeRequest(req)

	function := mtype.method.Func
	var returnValues []reflect.Value

	r := &request{
		service:     req.msg.Target,
		contentType: req.msg.Header["Content-Type"],
		method:      req.msg.Method,
		endpoint:    req.msg.Endpoint,
		body:        req.msg.Body,
	}

	// only set if not nil
	if argv.IsValid() {
		r.rawBody = argv.Interface()
	}

	if !mtype.stream {
		fn := func(ctx context.Context, req server.Request, rsp interface{}) error {
			returnValues = function.Call([]reflect.Value{s.rcvr, mtype.prepareContext(ctx), reflect.ValueOf(argv.Interface()), reflect.ValueOf(rsp)})

			// The return value for the method is an error.
			if err := returnValues[0].Interface(); err != nil {
				return err.(error)
			}

			return nil
		}

		// // wrap the handler
		// for i := len(router.hdlrWrappers); i > 0; i-- {
		// 	fn = router.hdlrWrappers[i-1](fn)
		// }

		// execute handler
		if err := fn(ctx, r, replyv.Interface()); err != nil {
			return err
		}

		// send response
		return router.sendResponse(sending, req, replyv.Interface(), cc, true)
	}

	return errors.New("Streaming Unsupported Now")
}

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// prepareMethod returns a methodType for the provided method or nil
// in case if the method was unsuitable.
func prepareMethod(method reflect.Method) *methodType {
	mtype := method.Type
	mname := method.Name
	var replyType, argType, contextType reflect.Type
	var stream bool

	// Method must be exported.
	if method.PkgPath != "" {
		return nil
	}

	switch mtype.NumIn() {
	case 3:
		// assuming streaming
		argType = mtype.In(2)
		contextType = mtype.In(1)
		stream = true
	case 4:
		// method that takes a context
		argType = mtype.In(2)
		replyType = mtype.In(3)
		contextType = mtype.In(1)
	default:
		log.Log("method", mname, "of", mtype, "has wrong number of ins:", mtype.NumIn())
		return nil
	}

	if !stream {
		// if not stream check the replyType
		// First arg need not be a pointer.
		if !isExportedOrBuiltinType(argType) {
			log.Log(mname, "argument type not exported:", argType)
			return nil
		}

		if replyType.Kind() != reflect.Ptr {
			log.Log("method", mname, "reply type not a pointer:", replyType)
			return nil
		}

		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			log.Log("method", mname, "reply type not exported:", replyType)
			return nil
		}
	}

	// Method needs one out.
	if mtype.NumOut() != 1 {
		log.Log("method", mname, "has wrong number of outs:", mtype.NumOut())
		return nil
	}
	// The return type of the method must be error.
	if returnType := mtype.Out(0); returnType != typeOfError {
		log.Log("method", mname, "returns", returnType.String(), "not error")
		return nil
	}
	return &methodType{method: method, ArgType: argType, ReplyType: replyType, ContextType: contextType, stream: stream}
}
