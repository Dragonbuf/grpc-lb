package configs

import (
	"flag"
	"sync"
)
import "github.com/Unknwon/goconfig"

var ConfigFilePath = flag.String("conf", "/data/code/go/grpc-lb/config-dev.ini", "config file path")
var once sync.Once
var cfg *goconfig.ConfigFile

func GetConfig() *goconfig.ConfigFile {

	once.Do(func() {
		flag.Parse()
		var err error
		cfg, err = goconfig.LoadConfigFile(*ConfigFilePath)
		if err != nil {
			panic(err)
		}
	})

	return cfg
}
