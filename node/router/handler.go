package router

import (
	"context"
	"encoding/xml"
	"regexp"

	_ "github.com/micro/go-micro/v2/broker"

	"github.com/micro/go-micro/v2/codec"
	"github.com/micro/go-micro/v2/util/log"

	eventbroker "github.com/micro-community/x-edge/broker"

	protocol "github.com/micro-community/x-edge/proto/protocol"
)

//ProtocolPackge means a private protocol pakcage
type ProtocolPackge struct {
	Version string `xml:"VER"`
	Name    string `xml:"NAME"`
	Gender  string `xml:"GENDER"`
	Type    string `xml:"TYPE"`
	Addr    string `xml:"ADDR"`
	Phone   string `xml:"PHONE"`
	Company string `xml:"COMPANY"`
	Time    string `xml:"TIME"`
}

type ProtocolServer struct {
}

func (e *ProtocolServer) unpackage(buf []byte) (protocolInfo ProtocolPackge, err error) {

	reg, _ := regexp.Compile("(?i:^<\\?xml (.+?)\\?>)")

	srcBuffer := buf
	indexs := reg.FindIndex(srcBuffer)
	if len(indexs) == 2 {
		srcBuffer = srcBuffer[indexs[1]+1:]
	}

	protocolInfo = ProtocolPackge{}
	err = xml.Unmarshal(srcBuffer, &protocolInfo)
	if err != nil {
		log.Errorf("unpackage error:%v", err)
		return protocolInfo, err
	}
	return protocolInfo, nil
}

func (e *ProtocolServer) createProtocMsg(protocolInfo ProtocolPackge) (msg protocol.Message) {
	msg = protocol.Message{
		Ver:     protocolInfo.Version,
		Name:    protocolInfo.Name,
		Gender:  protocolInfo.Gender,
		Type:    protocolInfo.Addr,
		Addr:    protocolInfo.Phone,
		Phone:   protocolInfo.Company,
		Company: protocolInfo.Time,
	}
	return msg
}

// Event is a single server request handler
func (e *ProtocolServer) Event(ctx context.Context, req *codec.Message, resp *codec.Message) error {

	logstr := string(req.Body[:])
	log.Log("[Received Protocol Event request]:", logstr)

	//unpackage from xml to ProtocolPackge
	ProtocolPackge, err := e.unpackage(req.Body)
	if err != nil {
		return nil
	}

	//rpc proto.NewProtocolService()

	//make the protocol.message
	buf := e.createProtocMsg(ProtocolPackge)
	//publish protocol.message
	err = eventbroker.PublishEventMessage(&buf)

	if err != nil {
		log.Errorf("eventbroker.Event PublishEventMessage error:%v", err)
	}
	return nil
}
