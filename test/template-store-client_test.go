package main

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"grpc-lb/api/protobuf-spec/template"
	"grpc-lb/internal/template/service"
	hystrix2 "grpc-lb/pkg/hystrix"
	"testing"
	"time"
)

func TestCanGetTemplate(t *testing.T) {
	client := service.NewTemplateClient()
	ctxc, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err := client.Get(ctxc, &template.ShowRequest{TemplateId: "T_08504EOI"})
	if err != nil {
		t.Error(err)
	}
}

func TestCanGetTemplateWithHystrix(t *testing.T) {

	_ = hystrix.Do(hystrix2.Name, func() error {
		client := service.NewTemplateClient()
		ctxc, _ := context.WithTimeout(context.Background(), 2*time.Second)
		_, err := client.Get(ctxc, &template.ShowRequest{TemplateId: "T_08504EOI"})
		return err
	}, nil)
}

func BenchmarkGetTemplate(n *testing.B) {
	for i := 0; i < n.N; i++ {
		client := service.NewTemplateClient()
		ctxc, _ := context.WithTimeout(context.Background(), 2*time.Second)
		_, err := client.Get(ctxc, &template.ShowRequest{TemplateId: "T_08504EOI"})
		if err != nil {
			panic(err)
		}
	}
}
