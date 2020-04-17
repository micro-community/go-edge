package transport

import (
	"bufio"
)

//DataExtractor for package pasering
type DataExtractor = bufio.SplitFunc

//DataExtractorFuncKey for DataExtractor
type DataExtractorFuncKey struct{}

var minDataPackageLength = 50

//DefaultdataExtractor will returns all data
func DefaultdataExtractor(data []byte, atEOF bool) (advance int, token []byte, err error) {

	if atEOF || len(data) == 0 {
		return 0, nil, nil
	}
	return len(data), data, nil
}
