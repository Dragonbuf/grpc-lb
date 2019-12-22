package main

import (
	"context"
	"grpc-lb/api/protobuf-spec/template"
	"grpc-lb/internal/common/loadBalance"
	"log"
	"testing"
)

func BenchmarkUploadConfig(n *testing.B) {

	roundrobinConn, err := loadBalance.NewBaseClient("template_store_service").GetRoundRobinConn()
	if err != nil {
		n.Error(err)
	}

	for i := 0; i < n.N; i++ {
		client := template.NewTemplateStoreClient(roundrobinConn)
		resp, err := client.Get(context.Background(), &template.ShowRequest{TemplateId: "T_3TYERB8E"})
		if err == nil && resp.TemplateId == "" {
			log.Fatal("error ")
		}
	}

}
