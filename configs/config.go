package configs

import (
	"flag"
)
import "github.com/Unknwon/goconfig"

var cfg *goconfig.ConfigFile
var ConfigFilePath = flag.String("conf", "/data/code/go/grpc-lb/config-dev.ini", "config file path")

func init() {
	flag.Parse()

	var err error
	cfg, err = goconfig.LoadConfigFile(*ConfigFilePath)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *goconfig.ConfigFile {
	return cfg
}
