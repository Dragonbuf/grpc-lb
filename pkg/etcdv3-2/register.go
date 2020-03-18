package etcdv3_2

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

type Config struct {
	ETCDEndpoints []string
}

var Deregister = make(chan struct{})

func NewClient(c *Config) *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: c.ETCDEndpoints,
	})
	if err != nil {
		panic(err)
	}
	return cli
}

// Register
func Register(etcdEndpoints []string, service, hostWithPort string, ttl int) error {

	serviceKey := fmt.Sprintf("/%s/%s/%s", schema, service, hostWithPort)

	// get endpoints for register dial address
	var err error
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: etcdEndpoints,
	})
	if err != nil {
		return fmt.Errorf("grpclb: create clientv3 client failed: %v", err)
	}
	resp, err := cli.Grant(context.TODO(), int64(ttl))
	if err != nil {
		return fmt.Errorf("grpclb: create clientv3 lease failed: %v", err)
	}

	if _, err := cli.Put(context.TODO(), serviceKey, hostWithPort, clientv3.WithLease(resp.ID)); err != nil {
		return fmt.Errorf("grpclb: set service '%s' with ttl to clientv3 failed: %s", service, err.Error())
	}

	if _, err := cli.KeepAlive(context.TODO(), resp.ID); err != nil {
		return fmt.Errorf("grpclb: refresh service '%s' with ttl to clientv3 failed: %s", service, err.Error())
	}

	// wait deregister then delete
	go func() {
		<-Deregister
		_, _ = cli.Delete(context.Background(), serviceKey)
		Deregister <- struct{}{}
	}()

	return nil
}

// UnRegister delete registered service from etcd
func UnRegister() {
	Deregister <- struct{}{}
	<-Deregister
}
