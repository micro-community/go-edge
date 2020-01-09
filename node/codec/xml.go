// Package codec provides a xml codec for default
package codec

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/micro/go-micro/codec"
)

//DefaultCodecs default Codec
var (
	DefaultCodecs = map[string]codec.NewCodec{
		"application/xml": NewCodec,
	}

	//DefaultContentType xml
	DefaultContentType = "application/xml"
)

//PROTOCOL Message Package Type
const (
	Error codec.MessageType = iota
	Request
	Response
	Event
	Notice
	File
)

//XMLBasicPackge represt a Protocol Frame
type XMLBasicPackge struct {
	Version string `xml:"VER"`
	Name    string `xml:"NAME"`
	Gender  string `xml:"GENDER"`
	Type    string `xml:"Type"`
}

//Codec for xml
type Codec struct {
	Conn    io.ReadWriteCloser
	Encoder *xml.Encoder
	Decoder *xml.Decoder
}

func finderRegString(regFormat string, srcString []byte) string {
	reg, _ := regexp.Compile(regFormat)

	indexs := reg.FindIndex(srcString)

	if indexs == nil || len(indexs) == 1 {
		return "nil"
	}

	targetString := string(srcString[indexs[0]:indexs[1]])

	return targetString
}

//ReadHeader return nil ,because no need of head
func (c *Codec) ReadHeader(m *codec.Message, t codec.MessageType) error {
	if m == nil || m.Body == nil {
		return nil
	}

	reg, _ := regexp.Compile("(?i:^<\\?xml (.+?)\\?>)")

	srcBuffer := m.Body

	indexs := reg.FindIndex(srcBuffer)

	if len(indexs) == 2 {
		m.Header["HeadLine"] = string(srcBuffer[0:indexs[1]])
		srcBuffer = srcBuffer[indexs[1]+1:]
	}

	basicInfo := XMLBasicPackge{}

	err := Marshaler{}.Unmarshal(srcBuffer, &basicInfo)

	if err != nil {
		return errors.New("Unmarshal XMLBasicPackge Struct error")
	}

	m.Header["VER"] = basicInfo.Version
	m.Header["NAME"] = basicInfo.Name
	m.Header["GENDER"] = basicInfo.Gender

	m.Target = "ProtocolServer"
	m.Endpoint = "protocol/" + basicInfo.Type
	m.Method = basicInfo.Type
	m.Body = srcBuffer

	if strings.EqualFold("File", basicInfo.Type) {
		m.Header["Stream"] = "true"
	}

	return nil
}

//ReadBody Get Data to handle
func (c *Codec) ReadBody(b interface{}) error {
	if b == nil {
		return nil
	}

	buf, err := ioutil.ReadAll(c.Conn)
	if err != nil {
		return err
	}

	reg, _ := regexp.Compile("(?i:^<\\?xml (.+?)\\?>)")

	srcBuffer := buf

	indexs := reg.FindIndex(srcBuffer)

	if len(indexs) == 2 {
		srcBuffer = srcBuffer[indexs[1]+1:]
	}

	return xml.Unmarshal(srcBuffer, b)

}

func (c *Codec) Write(m *codec.Message, b interface{}) error {
	if b == nil {
		return nil
	}
	return c.Encoder.Encode(b)
}

//Close stream
func (c *Codec) Close() error {
	return c.Conn.Close()
}

func (c *Codec) String() string {
	return "xml"
}

//NewCodec return xml codec
func NewCodec(c io.ReadWriteCloser) codec.Codec {
	return &Codec{
		Conn:    c,
		Decoder: xml.NewDecoder(c),
		Encoder: xml.NewEncoder(c),
	}
}
