package service

import (
	"context"
	"errors"
	"fmt"
	"grpc-lb/api/protobuf-spec/template"
	"grpc-lb/internal/common/loadBalance"
	model3 "grpc-lb/internal/templateStore/model"
	"sync"
)

var templateClient template.TemplateStoreClient
var once sync.Once

type server struct {
}

func NewTemplateServer() *server {
	return &server{}
}

func (t *server) Get(ctx context.Context, in *template.ShowRequest) (*template.ShowReply, error) {

	if in.GetTemplateId() == "" {
		return nil, errors.New("templateIdEmpty")
	}

	fmt.Println("templateId: " + in.TemplateId)
	model := model3.NewTemplateStoreModel()
	err := model.Get(in.GetTemplateId())
	if err != nil {
		return nil, err
	}

	return &template.ShowReply{
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

func NewTemplateClient() template.TemplateStoreClient {
	once.Do(func() {
		if templateClient == nil {
			fmt.Println("template client is nil, init ...")
			roundrobinConn, err := loadBalance.NewBaseClient("template_store_service").GetRoundRobinConn()
			if err != nil {
				panic(err)
			}
			templateClient = template.NewTemplateStoreClient(roundrobinConn)
		}
		fmt.Println("not do init again")
	})
	fmt.Println("already init template client")
	return templateClient
}
