package service

import (
	"context"
	"errors"
	"grpc-lb/api/protobuf-spec/template"
	"grpc-lb/pkg/loadBalance"
	"sync"
)

const Name = "template_store_service"

var templateClient template.TemplateClient
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

	return &template.ShowReply{TemplateId: in.TemplateId}, nil

	//fmt.Println("templateId: " + in.TemplateId)
	//model := model3.NewTemplateStoreModel()
	//err := model.Get(in.GetTemplateId())
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &template.ShowReply{TemplateId: model.TemplateId}, nil
}

func NewTemplateClient() template.TemplateClient {
	once.Do(func() {
		if templateClient == nil {
			roundRobinConn, err := loadBalance.NewBaseClient(Name).GetRoundRobinConn()
			if err != nil {
				panic(err)
			}
			templateClient = template.NewTemplateClient(roundRobinConn)
		}
	})
	return templateClient
}
