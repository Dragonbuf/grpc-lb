package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"grpc-lb/pkg/cache"
	"grpc-lb/pkg/db"
)

var configPath = ""
var Conf = &Config{}

func init() {
	flag.StringVar(&configPath, "config", "", "config path")
}

type Config struct {
	Redis *cache.Config
	Mysql *db.Config
}

func InitConf() error {
	if configPath != "" {
		if _, err := toml.DecodeFile(configPath, Conf); err != nil {
			return err
		}
	}
	return nil
}
