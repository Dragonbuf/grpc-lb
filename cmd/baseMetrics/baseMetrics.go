package baseMetrics

import (
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	_ "grpc-lb/cmd/config"
)

type InitMetrics struct {
	Reg *prometheus.Registry
	GrpcMetrics *grpc_prometheus.ServerMetrics
	//CustomizedCounterMetric *prometheus.CounterVec
}

func NewMetrics(name, help, label string) *InitMetrics {
	init := &InitMetrics{
		prometheus.NewRegistry(),
		grpc_prometheus.NewServerMetrics(),
		//prometheus.NewCounterVec(prometheus.CounterOpts{
		//	Name: name,
		//	Help: help,
		//}, []string{"naming"}),
	}
	init.Reg.MustRegister(init.GrpcMetrics)//,init.CustomizedCounterMetric
	//init.CustomizedCounterMetric.WithLabelValues(label)
	return init
}

func (b *InitMetrics) GetGpcServer() *grpc.Server {

	return grpc.NewServer(
		grpc.StreamInterceptor(b.GrpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(b.GrpcMetrics.UnaryServerInterceptor()),
	)
}
