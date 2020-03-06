package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	protocol "github.com/micro-community/x-edge/proto/protocol"
)

type Protocol struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Protocol) Call(ctx context.Context, req *protocol.Request, rsp *protocol.Response) error {
	log.Info("Received protocol.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Protocol) Stream(ctx context.Context, req *protocol.StreamingRequest, stream protocol.Protocol_StreamStream) error {
	log.Infof("Received Fun.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&protocol.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Protocol) PingPong(ctx context.Context, stream protocol.Protocol_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&protocol.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
