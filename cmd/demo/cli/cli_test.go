package main

import (
	"context"
	"flag"
	"testing"
	"time"

	"google.golang.org/grpc"

	_ "net/http/pprof"
)

func BenchmarkSend(c *testing.B) {
	flag.Parse()

	r := grpclb.NewResolver(*helloworld.serv)

	b := grpc.RoundRobin(r)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	conn, err := grpc.DialContext(ctx, *helloworld.reg, grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithBlock())
	cancel()
	if err != nil {
		panic(err)
	}

	for i := 0; i < c.N; i++ {
		client := pb.NewGreeterClient(conn)
		_, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "world "})
		if err == nil {
			//fmt.Printf("%v: Reply is %s\n", i, resp.Message)
		}
	}

}
