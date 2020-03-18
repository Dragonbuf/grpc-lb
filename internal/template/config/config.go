package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"grpc-lb/pkg/cache"
	"grpc-lb/pkg/db"
	etcdv3_2 "grpc-lb/pkg/etcdv3-2"
)

var configPath = ""
var Conf = &Config{}

func init() {
	flag.StringVar(&configPath, "config", "./test.toml", "config path")
}

type Server struct {
	Host string
	Port string
	Ttl  int
	Addr string
}

type Config struct {
	Redis  *cache.Config
	Mysql  *db.Config
	Etcd   *etcdv3_2.Config
	Server *Server
}

func InitConf() error {
	if configPath != "" {
		if _, err := toml.DecodeFile(configPath, Conf); err != nil {
			return err
		}
	}
	return nil
}
