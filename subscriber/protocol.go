package subscriber

import (
	"context"

	"github.com/micro/go-micro/util/log"

	protocol "github.com/micro-community/x-edge/proto/protocol"
)

type Protocol struct{}

func (e *Protocol) Handle(ctx context.Context, msg *protocol.Message) error {
	log.Log("Handler Received message: ", msg)

	return nil
}

func Handler(ctx context.Context, msg *protocol.Message) error {
	log.Log("Function Received message: ", msg)

	return nil
}
