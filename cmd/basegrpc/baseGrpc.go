package basegrpc

import (
	"context"
	"flag"
	"fmt"
	etcd3 "github.com/coreos/etcd/clientv3"
	grpclb "github.com/wwcd/grpc-lb/etcdv3"
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
	serv = flag.String("service", "hello_service", "service name")
	host = flag.String("host", "localhost", "listening host")
	port = flag.String("port", "50001", "listening port")
	reg  = flag.String("reg", "http://192.168.0.102:2379", "register etcd address")
	//reg  = flag.String("reg", "http://39.105.90.215:2379", "register etcd address")
	prefix = "/etcd3_naming"
)

type InitGrpc struct {
	ServiceName string
}

//TODO auto find port
func (b *InitGrpc) NewBaseGrpc() net.Listener {
	flag.Parse()

	*host = getLocalIp()
	*port = getLocalPort()

	lis, err := net.Listen("tcp", net.JoinHostPort(*host, *port))
	if err != nil {
		panic(err)
	}

	err = grpclb.Register(b.ServiceName, *host, *port, *reg, time.Second*10, 15)
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

	log.Printf("starting hello service at %s", *port)

	return lis
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
	endpoint = append(endpoint, *reg)

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
