package main

import (
	"google.golang.org/grpc"
	"grpc-lb/api/protobuf-spec/template"
	"grpc-lb/internal/template/service"
	"grpc-lb/pkg/loadBalance"
)

func main() {
	s := grpc.NewServer()
	template.RegisterTemplateServer(s, service.NewTemplateServer())
	_ = s.Serve(loadBalance.NewServer(service.Name).ReturnNetListenerWithRegisterLB())
}
