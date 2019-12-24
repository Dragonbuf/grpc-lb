package configs

import (
	"flag"
)
import "github.com/Unknwon/goconfig"

var cfg *goconfig.ConfigFile
var configFilePath = flag.String("conf", "./configs/config.ini", "config path")

func init() {
	flag.Parse()
	var err error
	cfg, err = goconfig.LoadConfigFile(*configFilePath)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *goconfig.ConfigFile {
	return cfg
}
