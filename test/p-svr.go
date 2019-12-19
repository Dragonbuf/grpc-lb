package main

import (
	"fmt"
	"grpc-lb/internal/pkg/baseMetrics"
	"grpc-lb/internal/pkg/loadBalance"
	"log"
	"net/http"

	pb "github.com/grpc-ecosystem/go-grpc-prometheus/examples/grpc-server-with-prometheus/protobuf"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"
)

// DemoServiceServer defines a Server.
type DemoServiceServer struct {
	Metrics *baseMetrics.InitMetrics
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

	server := loadBalance.NewServer("font")
	lis := server.GetAliveServer()

	metircs := baseMetrics.NewMetrics("hello_method_handle_count", "total", "name")
	// 带有 拦截器的 grpc 服务
	grpcServer := metircs.GetGrpcServer()

	service := newDemoServer()
	service.Metrics = metircs
	//  注册服务
	pb.RegisterDemoServiceServer(grpcServer, service)

	// 初始化 普罗米修斯
	metircs.GrpcMetrics.InitializeMetrics(grpcServer)

	// Start your http server for prometheus.
	go func() {

		// Create a HTTP server for prometheus.
		httpServer := &http.Server{Handler: promhttp.HandlerFor(metircs.Reg, promhttp.HandlerOpts{}), Addr: "0.0.0.0:9092"}

		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	// Start your gRPC server.
	_ = grpcServer.Serve(lis)
}
