package etcd_web_watch

import (
	"context"
	"encoding/json"
	"fmt"
	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	_ "github.com/coreos/etcd/mvcc/mvccpb"
	"log"
	"net/http"
	_ "net/http"
	"strings"
	_ "strings"
)

var Prefix = "/etcd3_naming"

type ServiceList struct {
	Service      []*ServiceAddrs
	AliveService int
	TotalService int
	ServiceName  string
}

type ServiceAddrs struct {
	IsExpire    bool
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
				if v.IsExpire == true {
					v.IsExpire = false
				}
				svr.svrMap[svrName].AliveService++
				isHave = true
			}
		}

		// // 已经注册了此服务，但是没有此地址
		if !isHave {
			addr := &ServiceAddrs{false, svrAddr}
			svr.svrMap[svrName].ServiceName = svrName
			svr.svrMap[svrName].Service = append(svr.svrMap[svrName].Service, addr)
			svr.svrMap[svrName].AliveService++
			svr.svrMap[svrName].TotalService++
		}
	} else {

		addr := &ServiceAddrs{false, svrAddr}
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
		log.Fatal("not found this service:" + svrName)
		return
	}

	for _, v := range svr.svrMap[svrName].Service {
		if v.ServiceAddr == svrAddr && v.IsExpire == false {
			v.IsExpire = true
			svr.svrMap[svrName].AliveService--
		}
	}

}

var svrMap ServiceMap

func Start() {

	var endpoint []string
	endpoint = append(endpoint, "localhost:2379")

	client, err := etcd3.New(etcd3.Config{
		Endpoints: endpoint,
	})
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Get(context.Background(), Prefix, etcd3.WithPrefix())
	if err != nil {
		panic(err)
	}

	// 注册 已经上线的服务
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			server := strings.TrimRight(strings.TrimLeft(string(resp.Kvs[i].Key), Prefix), string(resp.Kvs[i].Value))
			svrMap.AddMap(server, string(resp.Kvs[i].Value))
		}
	}

	go func() {
		// watch 观察是否有服务上线或下线
		ctx, cancel := context.WithCancel(context.Background())
		rch := client.Watch(ctx, Prefix, etcd3.WithPrefix(), etcd3.WithPrevKV())
		defer cancel()
		for wresp := range rch {
			for _, ev := range wresp.Events {
				switch ev.Type {
				case mvccpb.PUT:
					server := strings.TrimRight(strings.TrimLeft(string(ev.Kv.Key), Prefix), string(ev.Kv.Value))
					fmt.Println(ev.Kv)
					svrMap.AddMap(server, string(ev.Kv.Value))
				case mvccpb.DELETE:
					fmt.Println(ev.PrevKv)
					server := strings.TrimRight(strings.TrimLeft(string(ev.PrevKv.Key), Prefix), string(ev.PrevKv.Value))
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

	fmt.Println("server run in 127.0.0.1:8000")
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	jsonD, _ := json.Marshal(svrMap.svrMap)

	fmt.Println(string(jsonD))
	fmt.Fprintln(w, string(jsonD))
}

func test() {
	//svrMap.AddMap("shit","127.0.0.1:9000")
	//fmt.Println(svrMap)
	//svrMap.AddMap("shit2","127.0.0.1:9002")
	//fmt.Println(svrMap)
	//svrMap.DelMap("shit2","127.0.0.1:9002")
}
