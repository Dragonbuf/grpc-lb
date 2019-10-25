package main

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	model3 "grpc-lb/internal/app/server/templateStore/model"
	"grpc-lb/internal/app/server/templateStore/proto"
	"grpc-lb/internal/pkg/baseServer"
)

func main() {
	var base = baseServer.NewServer("template_store_service")
	lis := base.GetAliveServer()

	s := grpc.NewServer()
	templateStore.RegisterTemplateStoreServer(s, &server{})

	base.RegisterServer("template_store_copy_service")
	_ = s.Serve(lis)
}

type server struct {
}

func (t *server) Get(ctx context.Context, in *templateStore.ShowRequest) (*templateStore.ShowReply, error) {

	if in.GetTemplateId() == "" {
		return nil, errors.New("templateIdEmpty")
	}

	model := model3.NewTemplateStoreModel()
	err := model.Get(in.GetTemplateId())
	if err != nil {
		return nil, err
	}

	return &templateStore.ShowReply{
		TemplateId:       model.TemplateId,
		DesignerId:       model.DesignerId,
		Price:            model.Price,
		Version:          model.Version,
		Star:             model.Star,
		Level:            model.Level,
		Sort:             model.Sort,
		IsVipFree:        model.IsVipFree,
		TemplateInfo:     model.TemplateInfo,
		TemplateProperty: model.TemplateProperty,
	}, nil
}
