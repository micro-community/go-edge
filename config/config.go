package config

import (
	"fmt"

	mconfig "github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
)

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

//XMicroEdgeServiceTransport is the current srv MicroService Name
var XMicroEdgeServiceTransport = "tcp"

//XMicroEdgeServiceAddr ip Addr
var XMicroEdgeServiceAddr = "192.168.1.198:6600"

//Database define our own Database Config
type Database struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
}

//Cache define our own Cache Config
type Cache struct {
	Database
}

//MicroSets define our own Micro Config
type MicroSets struct {
	MicroServerName       string `toml:"microservername"`
	MicroServerAddress    string `toml:"microserveraddress"`
	MicroRegisterTTL      int    `toml:"microregisterttl"`
	MicroRegisterInterval int    `toml:"microregisterinterval"`
}

//Config From filea
var (
	DBConfig    Database
	CacheConfig Cache
	MicroConfig MicroSets
)

func init() {

	// load the config from a file source
	if err := mconfig.Load(file.NewSource(file.WithPath("./config.toml"))); err != nil {
		fmt.Println(err)
	}

	// read a Micro ENVVar
	if err := mconfig.Get("micro").Scan(&MicroConfig); err != nil {
		fmt.Println(err)
	}

	// read a database host
	if err := mconfig.Get("hosts", "database").Scan(&DBConfig); err != nil {
		fmt.Println(err)
	}

	// read a cache
	if err := mconfig.Get("hosts", "cache").Scan(&CacheConfig); err != nil {
		fmt.Println(err)
	}

	//	ServiceName = MicroConfig.ServeName

}
