package edge

import (
	ns "github.com/micro-community/x-edge/node/server"
	service "github.com/micro/go-micro/service"
)

// Options  of edge node serivices

//WithExtractor edge message
func WithExtractor(de ns.DataExtractor) service.Option {
	return func(o *service.Options) {
		o.Server.Init(ns.Extractor(de))
	}
}
