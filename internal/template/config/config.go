package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"grpc-lb/pkg/cache"
)

var configPath = ""
var Conf = &Config{}

func init() {
	flag.StringVar(&configPath, "config", "", "config path")
}

type Config struct {
	Redis *cache.Config
}

func NewConf() error {
	if configPath != "" {
		if _, err := toml.DecodeFile(configPath, Conf); err != nil {
			return err
		}
	}
	return nil
}
