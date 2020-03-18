// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"
	"grpc-lb/internal/template/config"
	"grpc-lb/pkg/loadBalance"
)

// InitializeEvent 声明injector的函数签名
func InitService(name string, cf *config.Config) *loadBalance.Service {
	wire.Build(
		loadBalance.WithDb,
	)
	return loadBalance.NewService(name) //返回值没有实际意义，只需符合函数签名即可
}
