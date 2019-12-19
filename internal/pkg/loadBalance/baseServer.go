package loadBalance

import (
	"flag"
	"grpc-lb/configs"
	etcdv3V2 "grpc-lb/pkg/etcdv3-2"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const defaultTTL = 15

var (
	host = flag.String("h", "127.0.0.1", "listening host")
	port = flag.String("p", "9001", "listening port")
)

type InitServer struct {
	ServiceName      string
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

// 得到一个 tcp 连接, 同时注册到 load balance 上面
func (b *InitServer) ReturnNetListenerWithRegisterLB() net.Listener {
	flag.Parse()

	lis, err := net.Listen("tcp", net.JoinHostPort(*host, *port))
	if err != nil {
		panic(err)
	}

	b.Register(b.ServiceName, b.getInputHostWithPort(), defaultTTL)

	return lis
}

func (b *InitServer) Register(serviceName, hostWithPort string, ttl int) {
	err := etcdv3V2.Register(configs.ETCDEndpoints, serviceName, hostWithPort, ttl)
	if err != nil {
		panic(err)
	}

	log.Printf("\n %s starting %s at etcdv3 \n", b.ServiceName, hostWithPort)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		etcdv3V2.UnRegister()
		os.Exit(1)
	}()
}

func (b *InitServer) getInputHostWithPort() string {
	if *host == "" {
		*host = getLocalIp()
	}
	if *port == "" {
		panic("port can not empty")
	}

	return *host + ":" + *port
}

// 如果要注册第二个服务，则传递服务名称即可
func (b *InitServer) AppendServiceWithSameServer(serviceName string) {
	b.Register(serviceName, b.getInputHostWithPort(), defaultTTL)
}

func getLocalIp() string {
	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		log.Printf("Get local IP addr failed!!!")
	}
	for _, addr := range addrSlice {
		if inet, ok := addr.(*net.IPNet); ok && !inet.IP.IsLoopback() {
			if nil != inet.IP.To4() {
				return inet.IP.String()
			}
		}
	}
	return "localhost"
}
