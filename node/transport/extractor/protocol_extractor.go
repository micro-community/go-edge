package extractor

import (
	"bufio"
)

var call bufio.SplitFunc

func DataExtractor(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return call(data, atEOF)
}

// RegisterExtractorHandler for Data Extractor
func RegisterExtractorHandler(model bufio.SplitFunc) {
	call = model
}
