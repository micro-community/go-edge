package config

//Sevice configuration and register
const (
	//RegisterTTL Time
	RegisterTTL = 30
	//RegisterInterval Time
	RegisterInterval = 10
	//Service Name
	XMicroEdgeServiceName = "x-micro-edge"
	// Version is a built-time-injected variable.
	XMicroEdgeServiceVersion = "1.0.0.0"
	//EventPublisherName is pubisher topic name
	EventPublisherName = "x-micro-edge.pubevent"
	//EventSubscriberName is subscriber topic name
	EventSubscriberName = "x-micro-edge.subevent"
)

//ServiceName is the current srv MicroService Name
var XMicroEdgeServiceTransport = "tcp"

//XMicroEdge ip Addr
var XMicroEdgeServiceAddr = "192.168.1.198:6600"
