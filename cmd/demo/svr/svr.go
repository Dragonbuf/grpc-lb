package main

import (
	"flag"
	_ "github.com/grpc-ecosystem/go-grpc-prometheus"
	"grpc-lb/api/protobuf-spec/template"
	"grpc-lb/internal/template/service"
	"grpc-lb/pkg/baseService"
)

func main() {
	flag.Parse()
	s := baseService.NewBaseService(service.Name)
	template.RegisterTemplateServer(s.GetGrpcServer(), service.NewTemplateServer())
	s.StartAndServe()
}
