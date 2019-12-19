package etcdv3_2

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/resolver"
	"grpc-lb/configs"
)

const schema = "etcdv3_resolver"

var addrs = []string{"localhost:50051", "localhost:50052"}

type ResolverBuilder struct {
	target     []string
	service    string
	etcdClient *clientv3.Client
	cc         resolver.ClientConn
}

// NewResolver return resolver builder
// target example: []{http://127.0.0.1:2379,http://127.0.0.1:12379,http://127.0.0.1:22379}
// service is service name
func NewResolver(target []string, service string) resolver.Builder {
	return &ResolverBuilder{target: target, service: service}
}

// ResolveNow
func (r *ResolverBuilder) ResolveNow(rn resolver.ResolveNowOption) {
}

// Close
func (r *ResolverBuilder) Close() {
}

// Scheme return etcdv3 schema
func (r *ResolverBuilder) Scheme() string {
	return schema
}

func (r *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {

	client, err := etcd3.New(etcd3.Config{
		Endpoints: configs.ETCDEndpoints,
	})
	if err != nil {
		return nil, fmt.Errorf("grpclb: creat etcd3 client failed %v", err)
	}

	r.etcdClient = client
	r.cc = cc

	go r.Watcher(fmt.Sprintf("/%s/%s/", schema, r.service))

	return r, nil
}

func (r *ResolverBuilder) Watcher(prefix string) {
	fmt.Println("watch :", prefix)
	addrDict := make(map[string]resolver.Address)

	update := func() {
		addrList := make([]resolver.Address, 0, len(addrDict))
		for _, v := range addrDict {
			addrList = append(addrList, v)
		}
		fmt.Println(addrList)
		r.cc.NewAddress(addrList)
		//r.cc.UpdateState(resolver.State{Addresses:     addrList})
	}

	resp, err := r.etcdClient.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err == nil {
		for i := range resp.Kvs {
			addrDict[string(resp.Kvs[i].Value)] = resolver.Address{Addr: string(resp.Kvs[i].Value)}
		}
	} else {
		panic(err)
	}

	update()

	rch := r.etcdClient.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for n := range rch {
		for _, ev := range n.Events {
			switch ev.Type {
			case mvccpb.PUT:
				addrDict[string(ev.Kv.Key)] = resolver.Address{Addr: string(ev.Kv.Value)}
			case mvccpb.DELETE:
				delete(addrDict, string(ev.PrevKv.Key))
			}
		}
		update()
	}
}
