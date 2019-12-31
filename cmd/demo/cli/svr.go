package main

import (
	"flag"
	"fmt"
	"grpc-lb/internal/template/config"
	"grpc-lb/pkg/cache"
	db2 "grpc-lb/pkg/db"
	_ "net/http/pprof"
)

func main() {
	flag.Parse()
	if err := config.InitConf(); err != nil {
		panic(err)
	}

	redis := cache.NewRedisPool(config.Conf.Redis)
	db := db2.NewMysql(config.Conf.Mysql)

	fmt.Println(redis, db)

}
