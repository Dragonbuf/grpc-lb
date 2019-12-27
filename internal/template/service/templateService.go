package service

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"grpc-lb/api/protobuf-spec/template"
	"grpc-lb/pkg/loadBalance"
	"sync"
)

const Name = "template_store_service"

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

var templateClient template.TemplateClient

type NewClientOptions func(c *TemplateClient)

type Middle func(ctx context.Context, in interface{}) (interface{}, error)
type EndPoint func(ctx context.Context, in interface{}) (interface{}, error)

type TemplateClient struct {
	name          string
	cc            *grpc.ClientConn
	tc            template.TemplateClient
	beforeActions []Middle
	Endpoints     []EndPoint
}

func NewTemplateClient(opts ...NewClientOptions) *TemplateClient {
	once.Do(func() {
		if templateClient == nil {
			roundRobinConn, err := loadBalance.NewBaseClient(Name).GetRoundRobinConn()
			if err != nil {
				panic(err)
			}
			templateClient = template.NewTemplateClient(roundRobinConn)
		}
	})
	client := &TemplateClient{tc: templateClient}
	for _, opt := range opts {
		opt(client) //opt是个方法，入参是*Client，内部会修改client的值
	}
	return client
}

func WithWarp(action Middle) NewClientOptions {
	return func(c *TemplateClient) {
		c.beforeActions = append(c.beforeActions, action)
	}
}

func WithWarpEndPoint(req EndPoint) NewClientOptions {
	return func(c *TemplateClient) {
		c.Endpoints = append(c.Endpoints, req)
	}
}

func WithName(name string) NewClientOptions {
	return func(c *TemplateClient) {
		c.name = name
	}
}

func (c *TemplateClient) Get(ctx context.Context, in *template.ShowRequest) (*template.ShowReply, error) {
	if len(c.beforeActions) > 0 {
		for _, fn := range c.beforeActions {
			if req, err := fn(ctx, in); err == nil {
				if templateReq, ok := req.(*template.ShowRequest); ok {
					in = templateReq
				}
			}
		}
	}

	if len(c.Endpoints) > 0 {
		for _, fn := range c.Endpoints {
			if reply, err := fn(ctx, c.tc.Get); err != nil {
				return nil, err
			} else {
				if templateShowReply, ok := reply.(*template.ShowReply); ok {
					return templateShowReply, nil
				}
			}
		}
	}

	return c.tc.Get(ctx, in)
}
