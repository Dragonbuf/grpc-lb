package helloworld

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"grpc-lb/api/protobuf-spec/helloworld"
	"grpc-lb/config"
	"strconv"
	"time"

	grpclb "grpc-lb/pkg/etcdv3"
)

var (
	serv = flag.String("service", "hello_service", "service name")
	reg  = flag.String("reg", config.EtcDHost, "register etcd address")
)

func main() {
	flag.Parse()
	r := grpclb.NewResolver(*serv)

	b := grpc.RoundRobin(r)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithBlock())
	cancel()
	if err != nil {
		panic(err)
	}

	//ticker := time.NewTicker(1000 * time.Millisecond)
	//for t := range ticker.C {
	t := time.Now()
	client := helloworld.NewGreeterClient(conn)
	resp, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
	if err == nil {
		fmt.Printf("%v: Reply is %s\n", t, resp.Message)
	}
	fmt.Println(time.Since(t))
	//}
}
