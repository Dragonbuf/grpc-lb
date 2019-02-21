package basegrpc

import (
	"flag"
	grpclb "github.com/wwcd/grpc-lb/etcdv3"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	serv = flag.String("service", "hello_service", "service name")
	host = flag.String("host", "localhost", "listening host")
	port = flag.String("port", "50001", "listening port")
	reg  = flag.String("reg", "http://192.168.1.171:2379", "register etcd address")
	//reg  = flag.String("reg", "http://39.105.90.215:2379", "register etcd address")
)



type BaseGrpc struct {
	ServiceName string
}
//TODO auto find port
func (b *BaseGrpc)NewBaseGrpc() net.Listener    {
	flag.Parse()

	lis, err := net.Listen("tcp", net.JoinHostPort(getLocalIp(), *port))
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
