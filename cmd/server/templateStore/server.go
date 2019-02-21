package main

import (
	"context"
	"errors"
	"github.com/wwcd/grpc-lb/cmd/basegrpc"
	"google.golang.org/grpc"
	model2 "grpc-lb/cmd/server/templateStore/model"
	ts "grpc-lb/cmd/server/templateStore/proto"
)

func main() {
	var base = basegrpc.BaseGrpc{ServiceName: "templateStoreService"}
	lis := base.NewBaseGrpc()

	s := grpc.NewServer()
	ts.RegisterTemplateStoreServer(s, &server{})
	_ = s.Serve(lis)
}

type server struct {
}

func (t *server) Get(ctx context.Context, in *ts.ShowRequest) (*ts.ShowReply, error) {

	if in.GetTemplateId() == "" {
		return nil, errors.New("templateIdEmpty")
	}

	var model model2.TemplateStoreModel
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
