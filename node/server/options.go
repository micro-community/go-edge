package server

import (
	"bufio"
	"context"

	"github.com/micro/go-micro/v2/codec"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/transport"
)

//DataExtractor for package pasering
type DataExtractor = bufio.SplitFunc

//DataExtractorFuncKey for ExtractorFunc
type DataExtractorFuncKey struct{}

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

func newOptions(opt ...server.Option) server.Options {
	opts := server.Options{
		Codecs:   make(map[string]codec.NewCodec),
		Metadata: map[string]string{},
	}

	for _, o := range opt {
		o(&opts)
	}

	if opts.Transport == nil {
		opts.Transport = transport.DefaultTransport
	}

	if len(opts.Address) == 0 {
		opts.Address = server.DefaultAddress
	}

	if len(opts.Name) == 0 {
		opts.Name = server.DefaultName
	}

	if len(opts.Id) == 0 {
		opts.Id = server.DefaultId
	}
	if len(opts.Version) == 0 {
		opts.Version = server.DefaultVersion
	}

	return opts
}

// Extractor should be used to setup a extractor
func Extractor(dex DataExtractor) server.Option {
	return func(o *server.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, DataExtractorFuncKey{}, dex)
	}
}
