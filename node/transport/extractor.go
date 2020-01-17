package transport

import (
	"regexp"
)

//DataExtractor for package pasering
type DataExtractor func(data []byte, atEOF bool) (advance int, token []byte, err error)

type dataExtractorFuncKey struct{}

var minDataPakckageLenth = 50

//extract data pakcage
func dataExtractor(data []byte, atEOF bool) (advance int, token []byte, err error) {

	if atEOF || len(data) == 0 {
		return 0, nil, nil
	}

	reg, _ := regexp.Compile("(?i:</protocol>)")

	indexs := reg.FindIndex(data)

	if indexs == nil || indexs[0] <= minDataPakckageLenth {
		return -1, data, nil //errors.New("error to extract data from socket")
		//return
	}

	advance = indexs[1]
	token = data[0:indexs[1]]
	return

}
