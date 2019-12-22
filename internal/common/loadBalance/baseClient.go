package loadBalance

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"grpc-lb/configs"
	etcdv3V2 "grpc-lb/internal/common/etcdv3-2"
	"time"
)

func init() {
	resolver.Register(&etcdv3V2.ResolverBuilder{})
}

type BaseClient struct {
	serviceName string
}

func NewBaseClient(serviceName string) *BaseClient {
	return &BaseClient{serviceName: serviceName}
}

func (c *BaseClient) GetRoundRobinConn() (*grpc.ClientConn, error) {

	r := etcdv3V2.NewResolver(configs.ETCDEndpoints, c.serviceName)
	resolver.Register(r)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return grpc.DialContext(ctx,
		r.Scheme()+"://authority/"+c.serviceName,
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
		grpc.WithBlock())
}
