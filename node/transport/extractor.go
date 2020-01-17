package transport

//DataExtractor for package pasering
type DataExtractor func(data []byte, atEOF bool) (advance int, token []byte, err error)

//DataExtractorFuncKey for DataExtractor
type DataExtractorFuncKey struct{}

var minDataPakckageLenth = 50
