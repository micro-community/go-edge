package subscriber

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	protocol "github.com/micro-community/x-edge/proto/protocol"
)

type Protocol struct{}

func (e *Protocol) Handle(ctx context.Context, msg *protocol.Message) error {
	log.Info("Handler Received message: ", msg)

	return nil
}

func Handler(ctx context.Context, msg *protocol.Message) error {
	log.Info("Function Received message: ", msg)

	return nil
}
