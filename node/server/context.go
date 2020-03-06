package server

import (
	"context"
	"sync"

	"github.com/micro/go-micro/v2/server"
)

type serverKey struct{}

func wait(ctx context.Context) *sync.WaitGroup {
	if ctx == nil {
		return nil
	}
	wg, ok := ctx.Value("wait").(*sync.WaitGroup)
	if !ok {
		return nil
	}
	return wg
}

//FromContext ...
func FromContext(ctx context.Context) (server.Server, bool) {
	c, ok := ctx.Value(serverKey{}).(server.Server)
	return c, ok
}

//NewContext ...
func NewContext(ctx context.Context, s server.Server) context.Context {
	return context.WithValue(ctx, serverKey{}, s)
}
