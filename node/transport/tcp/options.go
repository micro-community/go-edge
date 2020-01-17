package tcp

import (
	"context"
	"github.com/micro/go-micro/transport"

	nts "github.com/micro-community/x-edge/node/transport"
)

type dataExtractorFunc struct{}

// WithExtractor should be used to setup a extractor
func WithExtractor(dex nts.DataExtractor) transport.Option {
	return func(o *transport.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, dataExtractorFunc{}, dex)
	}
}
