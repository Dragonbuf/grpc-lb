package baseService

import (
	_ "github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"grpc-lb/pkg/baseMetrics"
	"grpc-lb/pkg/loadBalance"
	"net"
)

type BaseServiceOptions func(s *BaseService)

type BaseService struct {
	name       string
	grpcServer *grpc.Server
	metrics    *baseMetrics.InitMetrics
	lis        net.Listener
}

func NewBaseService(name string, opts ...BaseServiceOptions) *BaseService {
	metrics := baseMetrics.NewBaseMetrics()

	s := &BaseService{
		name:       name,
		metrics:    metrics,
		grpcServer: metrics.GetGrpcServer(),
		lis:        loadBalance.NewServer(name).ReturnNetListenerWithRegisterLB()}

	for _, opts := range opts {
		opts(s)
	}

	return s
}

func (s *BaseService) GetGrpcServer() *grpc.Server {
	return s.metrics.GetGrpcServer()
}

func (s *BaseService) StartAndServe() {

	s.metrics.NewSummaryForm(prometheus.SummaryOpts{
		Namespace: "template",
		Subsystem: "api",
		Name:      "per_request_duration",
	}, []string{"template_api_duration"}).With(prometheus.Labels{"template_api_duration": "start_and_serve"}).Observe(123456)

	s.metrics.NewCounterForm(prometheus.CounterOpts{
		Namespace: "template",
		Subsystem: "api",
		Name:      "count_all_request",
	}, []string{"template_api_duration"}).WithLabelValues("count_all_request").Add(1)

	s.metrics.InitAndServe()
	_ = s.grpcServer.Serve(s.lis)
}
