package codec

import (
	"encoding/xml"
)

//Marshaler for xml
type Marshaler struct{}

//Marshal for xml
func (x Marshaler) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

//Unmarshal for xml
func (x Marshaler) Unmarshal(b []byte, v interface{}) error {
	return xml.Unmarshal(b, v)
}

func (x Marshaler) String() string {
	return "xml"
}
