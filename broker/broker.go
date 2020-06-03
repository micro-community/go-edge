package eventbroker

import (
	"context"

	config "github.com/micro-community/x-edge/cmd"
	"github.com/micro-community/x-edge/proto/protocol"
	"github.com/micro-community/x-edge/subscriber"
	"github.com/micro/go-micro/v2"
)

// PubSubBroker is an status publisher for the go-micro broker
type PubSubBroker struct {
	publisher micro.Publisher
}

var (
	eventWorker PubSubBroker
)

// RegisterMessagePublisher creates a new broker status publisher
func RegisterMessagePublisher(service micro.Service) {
	eventWorker.publisher = micro.NewPublisher(config.EventPublisherName, service.Client())
}

// RegisterMessageSubscriber creates a new broker status subscriber
func RegisterMessageSubscriber(service micro.Service) {
	// Register Struct as Subscriber
	micro.RegisterSubscriber(config.EventSubscriberName, service.Server(), new(subscriber.Protocol))
	// Register Function as Subscriber
	micro.RegisterSubscriber(config.EventSubscriberName, service.Server(), subscriber.Handler)
}

//PublishEventMessage publish  protocol.message
func PublishEventMessage(msg *protocol.Message) error {
	return eventWorker.publisher.Publish(context.Background(), msg)
}
