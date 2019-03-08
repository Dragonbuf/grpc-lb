package main

import (
	"fmt"
	"grpc-lb/cmd/baseMetrics"
	"grpc-lb/cmd/baseServer"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	pb "github.com/grpc-ecosystem/go-grpc-prometheus/examples/grpc-server-with-prometheus/protobuf"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"
)

// DemoServiceServer defines a Server.
type DemoServiceServer struct{}

func newDemoServer() *DemoServiceServer {
	return &DemoServiceServer{}
}

// SayHello implements a interface defined by protobuf.
func (s *DemoServiceServer) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	customizedCounterMetric.WithLabelValues(request.Name).Inc()
	return &pb.HelloResponse{Message: fmt.Sprintf("Hello %s", request.Name)}, nil
}

var (
	// Create a metrics registry.
	reg2 = prometheus.NewRegistry()

	// Create some standard server metrics.
	grpcMetrics = grpc_prometheus.NewServerMetrics()

	// Create a customized counter metric.
	customizedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "demo_server_say_hello_method_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
)

func init() {
	// Register standard server metrics and customized metrics to registry.
	reg2.MustRegister(grpcMetrics, customizedCounterMetric)
	customizedCounterMetric.WithLabelValues("Test")
}

// NOTE: Graceful shutdown is missing. Don't use this demo in your production setup.
func main() {

	server := baseServer.NewServer("font")
	lis := server.GetAliveServer()

	// 带有 拦截器的 grpc 服务
	metrics := baseMetrics.NewMetrics()
	grpcServer := metrics.NewBaseMetrics()

	//  注册服务
	pb.RegisterDemoServiceServer(grpcServer, &DemoServiceServer{})

	// 初始化 普罗米修斯
	grpcMetrics.InitializeMetrics(grpcServer)

	// Start your http server for prometheus.
	go func() {

		// Create a HTTP server for prometheus.
		httpServer := &http.Server{Handler: promhttp.HandlerFor(reg2, promhttp.HandlerOpts{}), Addr: "0.0.0.0:9092"}

		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	// Start your gRPC server.
	_ = grpcServer.Serve(lis)
}
