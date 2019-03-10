package baseServer

import (
	"context"
	"flag"
	"fmt"
	etcd3 "github.com/coreos/etcd/clientv3"
	grpclb "github.com/wwcd/grpc-lb/etcdv3"
	"grpc-lb/cmd/config"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	serv   = flag.String("service", "", "service name")
	host   = flag.String("host", "", "listening host")
	port   = flag.String("port", "", "listening port")
	etcd   = flag.String("etcd", config.EtcDHost, "register etcd address")
	prefix = "/etcd3_naming"
)

type InitServer struct {
	ServiceName string
	ServiceNameSlice []string
}

type ServiceInfo struct {
	Name string
	Host string
	Port string
}


func NewServer(serName string) *InitServer {
	return &InitServer{ServiceName: serName}
}

// 得到一个 tcp 连接, 主动获得端口和内网 ip 地址
func (b *InitServer) GetAliveServer() net.Listener {
	flag.Parse()

	if *host == "" {
		*host = getLocalIp()
	}

	if *port == "" {
		*port = getLocalPort()
	}

	lis, err := net.Listen("tcp", net.JoinHostPort(*host, *port))
	if err != nil {
		panic(err)
	}

	err = grpclb.Register(b.ServiceName, *host, *port, *etcd, time.Second*10, 15)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		grpclb.UnRegister()
		os.Exit(1)
	}()

	log.Printf("\n starting %s at %s in etcd %s \n", b.ServiceName, *port, *etcd)

	return lis
}

// 如果要注册第二个服务，则传递服务名称即可
func (b *InitServer)RegisterServer(name string)  {
	err := grpclb.Register(name, *host, *port, *etcd, time.Second*10, 15)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v', close service :%s", s, name)
		grpclb.UnRegister()
		os.Exit(1)
	}()

	log.Printf("\n starting %s at %s in etcd %s \n", name, *port, *etcd)
}



func getLocalIp() string {
	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		log.Printf("Get local IP addr failed!!!")
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}

func getLocalPort() string {
	var endpoint []string
	endpoint = append(endpoint, *etcd)

	client, err := etcd3.New(etcd3.Config{
		Endpoints: endpoint,
	})

	if err != nil {
		panic(err)
	}

	resp, err := client.Get(context.Background(), prefix, etcd3.WithPrefix())
	if err != nil {
		panic(err)
	}

	ip := *host
	port := 50001
	if resp == nil || resp.Kvs == nil {
		fmt.Println("port can be 50001")
		return "50001"
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			if strings.Contains(string(v), ip+":"+strconv.Itoa(port)) {
				port++
			}
		}
	}

	lis, err := net.Listen("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	fmt.Println("can be run in " + strconv.Itoa(port))
	return strconv.Itoa(port)
}
