package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	pb "github.com/grpc-ecosystem/go-grpc-prometheus/examples/grpc-server-with-prometheus/protobuf"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"grpc-lb/internal/pkg/baseServer"
	"net/http"
)

type DemoServiceServer struct {
}

func newDemoServer() *DemoServiceServer {
	return &DemoServiceServer{}
}

// SayHello implements a interface defined by protobuf.
func (s *DemoServiceServer) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	//s.Metrics.CustomizedCounterMetric.WithLabelValues(request.Name).Inc()
	return &pb.HelloResponse{Message: fmt.Sprintf("Hello %s", request.Name)}, nil
}

// NOTE: Graceful shutdown is missing. Don't use this demo in your production setup.
func main() {

	server := baseServer.NewServer("font")
	lis := server.GetAliveServer()

	myServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)

	myService := newDemoServer()
	pb.RegisterDemoServiceServer(myServer, myService)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9093", nil)
	}()

	// Start your gRPC server.
	_ = myServer.Serve(lis)
}
