package main

import (
	"context"
	ts "grpc-lb/internal/app/server/templateStore/proto"
	"grpc-lb/internal/pkg/loadBalance"
	"log"
	"testing"
)

func BenchmarkUploadConfig(n *testing.B) {

	roundrobinConn, err := loadBalance.NewBaseClient("template_store_service").GetRoundRobinConn()
	if err != nil {
		n.Error(err)
	}

	for i := 0; i < n.N; i++ {
		client := ts.NewTemplateStoreClient(roundrobinConn)
		resp, err := client.Get(context.Background(), &ts.ShowRequest{TemplateId: "T_3TYERB8E"})
		if err == nil && resp.TemplateId == "" {
			log.Fatal("error ")
		}
	}

}
