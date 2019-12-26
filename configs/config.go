package configs

import (
	"flag"
)
import "github.com/Unknwon/goconfig"

var cfg *goconfig.ConfigFile
var configFilePath = flag.String("conf", "/data/code/go/grpc-lb/config-dev.ini", "config path")

func init() {
	var err error
	cfg, err = goconfig.LoadConfigFile(*configFilePath)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *goconfig.ConfigFile {
	return cfg
}
