package server

import (
	"context"

	"github.com/micro/go-micro/codec"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"
)

//DataExtractor for package pasering
type DataExtractor func(data []byte, atEOF bool) (advance int, token []byte, err error)

type dataExtractorFunc struct{}

// WithExtractor should be used to setup a extractor
func WithExtractor(dex DataExtractor) server.Option {
	return func(o *server.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, dataExtractorFunc{}, dex)
	}
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
