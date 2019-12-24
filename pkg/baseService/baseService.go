package baseService

import (
	"google.golang.org/grpc"
	"grpc-lb/pkg/baseMetrics"
	"grpc-lb/pkg/loadBalance"
	"net"
)

type BaseService struct {
	name       string
	grpcServer *grpc.Server
	metrics    *baseMetrics.InitMetrics
	lis        net.Listener
}

func NewBaseService(name string) *BaseService {
	metrics := baseMetrics.NewBaseMetrics()
	return &BaseService{
		name:       name,
		metrics:    metrics,
		grpcServer: metrics.GetGrpcServer(),
		lis:        loadBalance.NewServer(name).ReturnNetListenerWithRegisterLB()}
}

func (s *BaseService) GetGrpcServer() *grpc.Server {
	return s.metrics.GetGrpcServer()
}

func (s *BaseService) StartAndServe() {
	s.metrics.InitAndServe()
	_ = s.grpcServer.Serve(s.lis)
}
