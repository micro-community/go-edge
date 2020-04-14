package router

import (
	"context"
	"encoding/xml"
	"regexp"

	"github.com/micro/go-micro/v2/codec"
	log "github.com/micro/go-micro/v2/logger"

	eventbroker "github.com/micro-community/x-edge/broker"

	protocol "github.com/micro-community/x-edge/proto/protocol"
)

//ProtocolPackage means a private protocol package
type ProtocolPackage struct {
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

func (e *ProtocolServer) unpackage(buf []byte) (protocolInfo ProtocolPackage, err error) {

	reg, _ := regexp.Compile("(?i:^<\\?xml (.+?)\\?>)")

	srcBuffer := buf
	indexs := reg.FindIndex(srcBuffer)
	if len(indexs) == 2 {
		srcBuffer = srcBuffer[indexs[1]+1:]
	}

	protocolInfo = ProtocolPackage{}
	err = xml.Unmarshal(srcBuffer, &protocolInfo)
	if err != nil {
		log.Errorf("unpackage error:%v", err)
		return protocolInfo, err
	}
	return protocolInfo, nil
}

func (e *ProtocolServer) createProtocMsg(protocolInfo ProtocolPackage) (msg protocol.Message) {
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
	log.Info("[Received Protocol Event request]:", logstr)

	//unpackage from xml to ProtocolPackage
	ProtocolPackage, err := e.unpackage(req.Body)
	if err != nil {
		return nil
	}

	//rpc proto.NewProtocolService()

	//make the protocol.message
	buf := e.createProtocMsg(ProtocolPackage)
	//publish protocol.message
	err = eventbroker.PublishEventMessage(&buf)

	if err != nil {
		log.Errorf("eventbroker.Event PublishEventMessage error:%v", err)
	}
	return nil
}
