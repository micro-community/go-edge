package udp

import (
	"context"

	nts "github.com/micro-community/x-edge/node/transport"
)


func fromContext(ctx context.Context) (nts.DataExtractor, bool) {
	e, ok := ctx.Value(nts.DataExtractorFuncKey{}).(nts.DataExtractor)
	return e, ok
}
