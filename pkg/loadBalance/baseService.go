package loadBalance

import (
	etcd3 "github.com/coreos/etcd/clientv3"
	redis2 "github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"grpc-lb/internal/template/config"
	etcdv3V2 "grpc-lb/pkg/etcdv3-2"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	name  string
	redis *redis2.Pool
	db    *gorm.DB
	etcd  *etcd3.Client
	conf  *config.Config
	s     *grpc.Server
	log   *zap.Logger
}

type Options func(s *Service)

func NewService(name string, opts ...Options) *Service {
	s := &Service{}
	s.name = name
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithRedis(r *redis2.Pool) Options {
	return func(s *Service) {
		s.redis = r
	}
}

func WithDb(d *gorm.DB) Options {
	return func(s *Service) {
		s.db = d
	}
}

func WithZapLog(log *zap.Logger) Options {
	return func(s *Service) {
		s.log = log
	}
}

func WithGrpcServer(srv *grpc.Server) Options {
	return func(s *Service) {
		s.s = srv
	}
}

func WithConfig(c *config.Config) Options {
	return func(s *Service) {
		s.conf = c
	}
}

func WithEtcd(e *etcd3.Client) Options {
	return func(s *Service) {
		s.etcd = e
	}
}

func (s *Service) Start() error {
	lis, err := net.Listen("tcp", net.JoinHostPort(s.conf.Server.Host, s.conf.Server.Port))
	if err != nil {
		return err
	}

	// 注册服务到 etcdv3 todo 是否单独分离出 etcdv3
	if err := etcdv3V2.Register(s.conf.Etcd.ETCDEndpoints, s.name, s.conf.Server.Addr, s.conf.Server.Ttl); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		etcdv3V2.UnRegister()
		os.Exit(1)
	}()

	if err := s.s.Serve(lis); err != nil {
		return err
	}

	s.log.Info("server is start : " + s.conf.Server.Addr)
	return nil
}
