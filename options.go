package edge

import (
	ns "github.com/micro-community/x-edge/node/server"
	mservice "github.com/micro/go-micro/service"
)

// Options  of edge node serivices

//WithExtractor edge message
func (e ns.DataExtractor) mservice.Option {
	return func(o *service.Options) {
		o.Server.Init(ns.WithExtractor(e))
	}
}
