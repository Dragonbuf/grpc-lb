package etcd_web_watch

import (
	"context"
	"encoding/json"
	"fmt"
	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	_ "github.com/coreos/etcd/mvcc/mvccpb"
	"grpc-lb/configs"
	log "grpc-lb/internal/common/log"
	"net/http"
	_ "net/http"
	_ "strings"
)

type ServiceList struct {
	Service      []*ServiceAddrs
	AliveService int
	TotalService int
	ServiceName  string
}

type ServiceAddrs struct {
	ServiceAddr string
}

type ServiceMap struct {
	svrMap map[string]*ServiceList
}

func (svr *ServiceMap) AddMap(svrName string, svrAddr string) {
	isHave := false
	// 已经注册了此服务，并且有此地址
	if svr.svrMap[svrName] != nil {
		for _, v := range svr.svrMap[svrName].Service {
			if v.ServiceAddr == svrAddr {
				svr.svrMap[svrName].AliveService++
				isHave = true
			}
		}

		// // 已经注册了此服务，但是没有此地址
		if !isHave {
			addr := &ServiceAddrs{svrAddr}
			svr.svrMap[svrName].ServiceName = svrName
			svr.svrMap[svrName].Service = append(svr.svrMap[svrName].Service, addr)
			svr.svrMap[svrName].AliveService++
			svr.svrMap[svrName].TotalService++
		}
	} else {

		addr := &ServiceAddrs{svrAddr}
		s := &ServiceList{}

		s.ServiceName = svrName
		s.Service = []*ServiceAddrs{addr}
		s.AliveService++
		s.TotalService++

		// 如果 map 中没有任何数据，创建此map
		if svr.svrMap == nil {
			svr.svrMap = make(map[string]*ServiceList)
		}

		svr.svrMap[svrName] = s
	}

}

func (svr *ServiceMap) DelMap(svrName string, svrAddr string) {

	if svr.svrMap[svrName] == nil {
		log.GetLogger().Info("not found this service:" + svrName)
		return
	}

	for k, v := range svr.svrMap[svrName].Service {
		if v.ServiceAddr == svrAddr {
			//v.IsExpire = true
			// 直接删除掉 此 service
			svr.svrMap[svrName].Service = append(svr.svrMap[svrName].Service[:k], svr.svrMap[svrName].Service[k+1:]...)
			svr.svrMap[svrName].AliveService--
		}
	}

}

var svrMap ServiceMap

func Start() {

	client, err := etcd3.New(etcd3.Config{
		Endpoints: configs.ETCDEndpoints,
	})
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Get(context.Background(), configs.EtcdPrefix, etcd3.WithPrefix())
	if err != nil {
		panic(err)
	}

	// 注册 已经上线的服务
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {

			left := string(resp.Kvs[i].Key)[len(configs.EtcdPrefix) : len(string(resp.Kvs[i].Key))-1]
			server := left[0 : len(left)-len(string(resp.Kvs[i].Value))]

			//server := strings.TrimRight(strings.TrimLeft(string(resp.Kvs[i].Key), Prefix), string(resp.Kvs[i].Value))
			svrMap.AddMap(server, string(resp.Kvs[i].Value))
		}
	}

	go func() {
		// watch 观察是否有服务上线或下线
		ctx, cancel := context.WithCancel(context.Background())
		rch := client.Watch(ctx, configs.EtcdPrefix, etcd3.WithPrefix(), etcd3.WithPrevKV())
		defer cancel()
		for wresp := range rch {
			for _, ev := range wresp.Events {
				switch ev.Type {
				case mvccpb.PUT:
					left := string(ev.Kv.Key)[len(configs.EtcdPrefix) : len(string(ev.Kv.Key))-1]
					server := left[0 : len(left)-len(string(ev.Kv.Value))]
					svrMap.AddMap(server, string(ev.Kv.Value))
				case mvccpb.DELETE:
					left := string(ev.PrevKv.Key)[len(configs.EtcdPrefix) : len(string(ev.PrevKv.Key))-1]
					server := left[0 : len(left)-len(string(ev.PrevKv.Value))]
					svrMap.DelMap(server, string(ev.PrevKv.Value))
				}
			}
		}
	}()

	// http 展示服务
	http.HandleFunc("/", IndexHandler)

	err = http.ListenAndServe("127.0.0.1:8000", nil)
	if err != nil {
		panic(err)
	}

	log.GetLogger().Info("server run in 127.0.0.1:8000")
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	jsonD, _ := json.Marshal(svrMap.svrMap)

	fmt.Println(string(jsonD))
	_, _ = fmt.Fprintln(w, string(jsonD))
}
