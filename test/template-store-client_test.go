package main

import (
	"context"
	"flag"
	grpclb "github.com/wwcd/grpc-lb/etcdv3"
	"google.golang.org/grpc"
	ts "grpc-lb/internal/app/server/templateStore/proto"
	"log"
	"testing"
	"time"
)

func BenchmarkUploadConfig(n *testing.B) {
	flag.Parse()
	r := grpclb.NewResolver(*serv)

	b := grpc.RoundRobin(r)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithBlock())
	cancel()
	if err != nil {
		panic(err)
	}

	for i := 0; i < n.N; i++ {
		client := ts.NewTemplateStoreClient(conn)

		resp, err := client.Get(context.Background(), &ts.ShowRequest{TemplateId: "T_3TYERB8E"})
		if err == nil && resp.TemplateId == "" {
			log.Fatal("error ")
		}
	}

}
