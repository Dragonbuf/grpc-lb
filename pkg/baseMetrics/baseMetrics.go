package baseMetrics

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"grpc-lb/internal/common/log"
	"grpc-lb/internal/template/service"
	"net/http"
)

var (
	reg = prometheus.NewRegistry()

	// Create some standard server metrics.
	grpcMetrics = grpc_prometheus.NewServerMetrics()

	// Create a customized counter metric.
	customizedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "server_template_get_method_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{service.Name})

	sentBytes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "etcd",
		Subsystem: "network",
		Name:      "client_grpc_sent_bytes_total",
		Help:      "THe total number of bytes send to grpc clients",
	})

	receivedBytes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "etcd",
		Subsystem: "network",
		Name:      "client_grpc_received_bytes_total",
		Help:      "THe total number of bytes received from  grpc clients",
	})
)

func init() {
	// Create a metrics registry.
	reg.MustRegister(grpcMetrics, customizedCounterMetric, sentBytes, receivedBytes)
	customizedCounterMetric.WithLabelValues(service.Name)
	// Create a HTTP server for prometheus.

}

type InitMetrics struct {
	Reg         *prometheus.Registry
	GrpcMetrics *grpc_prometheus.ServerMetrics
	GrpcServer  *grpc.Server
}

func NewBaseMetrics() *InitMetrics {
	return &InitMetrics{
		GrpcServer: grpc.NewServer(
			grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
			grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
		),
	}
}

func (i *InitMetrics) GetGrpcServer() *grpc.Server {
	return i.GrpcServer
}

func (i *InitMetrics) InitAndServe() {
	grpcMetrics.InitializeMetrics(i.GetGrpcServer())

	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 9092)}

	// Start http server for prometheus.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.GetLogger().Error(err)
		}
	}()
}
