package main

import (
	"flag"
	_ "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	config "grpc-lb/configs"
	etcdv3V2 "grpc-lb/internal/common/etcdv3-2"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

var (
	serv = flag.String("service", "hello_service", "service name")
	host = flag.String("host", "localhost", "listening host")
	port = flag.String("port", "50001", "listening port")
	reg  = flag.String("reg", config.EtcDHost, "register etcd address")
	//reg  = flag.String("reg", "http://39.105.90.215:2379", "register etcd address")
)

func main() {
	flag.Parse()
	//*host = getLocalIp()
	lis, err := net.Listen("tcp", net.JoinHostPort(*host, *port))
	if err != nil {
		panic(err)
	}

	hostWithPort := *host + ":" + *port
	err = etcdv3V2.Register(config.ETCDEndpoints, *serv, hostWithPort, 15)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		etcdv3V2.UnRegister()
		os.Exit(1)
	}()

	log.Printf("starting hello service at %s", *port)

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	// Register your gRPC service implementations.
	//pb.RegisterGreeterServer(s, &server{})
	// After all your registrations, make sure all of the Prometheus metrics are initialized.
	grpc_prometheus.Register(s)

	// Register Prometheus metrics handler.
	http.Handle("/metrics", promhttp.Handler())

	//s := grpc.NewServer()
	//pb.RegisterGreeterServer(s, &server{})
	go func() {
		_ = http.ListenAndServe("localhost:13001", nil)
	}()
	_ = s.Serve(lis)

}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
//func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
//	fmt.Printf("%v: Receive is %s\n", time.Now(), in.Name)
//	return &pb.HelloReply{Message: "Hello " + in.Name + " from " + net.JoinHostPort(*host, *port)}, nil
//}
