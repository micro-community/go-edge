package edge

import (
	nts "github.com/micro-community/x-edge/node/transport"
	service "github.com/micro/go-micro/service"
)

// Options  of edge node serivices

//WithExtractor edge message
func WithExtractor(de nts.DataExtractor) service.Option {
	return func(o *service.Options) {
		o.Transport.Init(nts.WithExtractor(de))
	}
}
