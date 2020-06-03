package cmd

import (
	"fmt"

	mconfig "github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
)

//Service configuration and register
const (
	//RegisterTTL Time
	RegisterTTL = 30
	//RegisterInterval Time
	RegisterInterval = 10
	//EventPublisherName is pubisher topic name
	EventPublisherName = "x.edge.pubevent"
	//EventSubscriberName is subscriber topic name
	EventSubscriberName = "x.edge.subevent"
	//	Transport     = "udp"
	Host = ":8080"
	//  AppNamespace = "x.edge"
	HeaderPrefix = "x-edge-"
)

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
	ServerName       string `toml:"microservername"`
	ServerAddress    string `toml:"microserveraddress"`
	RegisterTTL      int    `toml:"microregisterttl"`
	RegisterInterval int    `toml:"microregisterinterval"`
}

//Config From files
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
