package main

import (
	"google.golang.org/grpc"
	"grpc-lb/api/protobuf-spec/template"
	"grpc-lb/internal/common/loadBalance"
	"grpc-lb/internal/templateStore/service"
)

func main() {
	srv := loadBalance.NewServer("template_store_service")

	s := grpc.NewServer()
	template.RegisterTemplateStoreServer(s, service.NewTemplateServer())
	_ = s.Serve(srv.ReturnNetListenerWithRegisterLB())
}
