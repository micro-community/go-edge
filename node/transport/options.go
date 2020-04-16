package transport

import (
	"context"

	"github.com/micro/go-micro/v2/transport"
)

// WithExtractor should be used to setup a extractor
func WithExtractor(dex DataExtractor) transport.Option {
	return func(o *transport.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, DataExtractorFuncKey{}, dex)
	}
}
