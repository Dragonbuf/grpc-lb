package main

import (
	"flag"
	"fmt"
	"grpc-lb/internal/template/config"
	"grpc-lb/pkg/cache"
	"testing"
)

func TestConfig(t *testing.T) {
	flag.Parse()
	if err := config.NewConf(); err != nil {
		t.Error(err)
	}
	fmt.Println(config.Conf.Redis)
	redis := cache.NewRedisPool(config.Conf.Redis)
	conn := redis.Get()
	defer conn.Close()
}
