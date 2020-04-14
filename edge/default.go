package edge

import (
	"errors"
	//here should edge internal logic
	"github.com/google/uuid"

	nclient "github.com/micro-community/x-edge/node/client"
	nserver "github.com/micro-community/x-edge/node/server"
	"github.com/micro-community/x-edge/node/transport/tcp"
	"github.com/micro-community/x-edge/node/transport/udp"
	"github.com/micro/go-micro/v2/transport"
)

// default components for edge node server
var (
	DefaultName           = "x-edge-node"
	DefaultVersion        = "latest"
	DefaultAddress        = ":8000"
	DefaultID             = uuid.New().String()
	DefaultClient         = nclient.NewClient()
	DefaultServer         = nserver.NewServer()
	DefaultTransport      = tcp.NewTransport()
	ErrNoExtractorDefined = errors.New("No Extractor Defined")

	DefaultExtractor = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return -1, nil, ErrNoExtractorDefined
	}
	DefaultTransports = map[string]func(...transport.Option) transport.Transport{
		"udp": udp.NewTransport,
		"tcp": tcp.NewTransport,
	}
)
