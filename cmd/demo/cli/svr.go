package main

import (
	"flag"
	zap2 "go.uber.org/zap"
	"google.golang.org/grpc"
	"grpc-lb/internal/template/config"
	"grpc-lb/pkg/cache"
	db2 "grpc-lb/pkg/db"
	etcdv3_2 "grpc-lb/pkg/etcdv3-2"
	"grpc-lb/pkg/loadBalance"
	_ "net/http/pprof"
)

func main() {
	flag.Parse()
	if err := config.InitConf(); err != nil {
		panic(err)
	}

	// todo 里面的内容使用 wire 依赖注入
	redis := cache.NewRedisPool(config.Conf.Redis)
	db := db2.NewMysql(config.Conf.Mysql)
	etcd := etcdv3_2.NewClient(config.Conf.Etcd)
	srv := grpc.NewServer()
	// todo 根据 config 使用 zap 的不同模式
	zap, err := zap2.NewDevelopment()
	if err != nil {
		panic(err)
	}

	service := loadBalance.NewService(
		"test",
		loadBalance.WithDb(db),
		loadBalance.WithRedis(redis),
		loadBalance.WithGrpcServer(srv),
		loadBalance.WithConfig(config.Conf),
		loadBalance.WithZapLog(zap),
		loadBalance.WithEtcd(etcd))

	if err := service.Start(); err != nil {
		panic(err)
	}
}
