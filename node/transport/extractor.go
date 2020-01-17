package transport

import (
	"bufio"
)

//DataExtractor for package pasering
type DataExtractor = bufio.SplitFunc

//DataExtractorFuncKey for DataExtractor
type DataExtractorFuncKey struct{}

var minDataPakckageLenth = 50
