package tcp

import (
	"context"

	nts "github.com/micro-community/x-edge/node/transpot"
)

type dataExtractorFuncKey struct{}

func fromContext(ctx context.Context) (nts.DataExtractor, bool) {
	e, ok := ctx.Value(dataExtractorFuncKey{}).(nts.DataExtractor)
	return e, ok
}
