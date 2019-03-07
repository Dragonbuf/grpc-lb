package main

import (
	"context"
	"flag"
	"fmt"
	"grpc-lb/cmd/config"
	"grpc-lb/cmd/server/templateStore/proto"
	"time"

	grpclb "github.com/wwcd/grpc-lb/etcdv3"
	"google.golang.org/grpc"
	ts "grpc-lb/cmd/server/templateStore/proto"
)

var (
	serv = flag.String("service", "template_store_service", "service name")
	reg  = flag.String("reg", config.EtcDHost, "register etcd address")
)

func main2() {
	flag.Parse()
	r := grpclb.NewResolver(*serv)

	b := grpc.RoundRobin(r)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithBlock())
	cancel()
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(1000 * time.Millisecond)
	for range ticker.C {
		client := templateStore.NewTemplateStoreClient(conn)

		resp, err := client.Get(context.Background(), &ts.ShowRequest{TemplateId: "T_3TYERB8E"})
		if err == nil {
			fmt.Printf("%v: \n", resp)
		}
	}
}
