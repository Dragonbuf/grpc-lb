package main

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"grpc-lb/cmd/baseServer"
	model2 "grpc-lb/cmd/server/templateStore/model"
	ts "grpc-lb/cmd/server/templateStore/proto"
)

func main() {
	var base = baseServer.NewServer("template_store_service")
	lis := base.GetAliveServer()

	s := grpc.NewServer()
	ts.RegisterTemplateStoreServer(s, &server{})


	base.RegisterServer("template_store_copy_service")
	_ = s.Serve(lis)
}

type server struct {
}

func (t *server) Get(ctx context.Context, in *ts.ShowRequest) (*ts.ShowReply, error) {

	if in.GetTemplateId() == "" {
		return nil, errors.New("templateIdEmpty")
	}

	model := model2.NewTemplateStoreModel()
	err := model.Get(in.GetTemplateId())
	if err != nil {
		return nil, err
	}

	return &ts.ShowReply{
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
