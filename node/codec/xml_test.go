package codec

import (
	"testing"

	"github.com/micro/go-micro/v2/codec"
)

var srcbytes = []byte(`<?xml version="1.0" encoding="gb2312"?>
<PROTOCOL>
<VER>1.0</VER>
<NAME>danny</NAME>
<GENDER>MALE</GENDER>
<TYPE>1</TYPE>
<ADDR>Road.1</ADDR>
<PHONE>400-800-5555</PHONE>
<COMPANY>xxx</COMPANY>
<TIME>2019.12.1-11:11:11</TIME>
</PROTOCOL>
`)


func TestReadHeader(t *testing.T) {

	msg := &codec.Message{Body: srcbytes}
	msg.Type = codec.Request
	msg.Header = map[string]string{}

	cdc := Codec{}
	cdc.ReadHeader(msg, msg.Type)

	t.Log("Done")

}
