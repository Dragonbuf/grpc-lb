package baseMetrics

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"grpc-lb/internal/common/log"
	_ "grpc-lb/internal/template/service"
	"net/http"
	"time"
)

var (
	reg = prometheus.NewRegistry()

	// Create some standard server metrics.
	grpcMetrics = grpc_prometheus.NewServerMetrics()
)

func init() {
	// Create a metrics registry.
	//reg.MustRegister(grpcMetrics, customizedCounterMetric, requestDuration)
	//customizedCounterMetric.WithLabelValues(service.Name).Add(1)
	// Create a HTTP server for prometheus.
}

type InitMetrics struct {
	Reg         *prometheus.Registry
	GrpcMetrics *grpc_prometheus.ServerMetrics
	GrpcServer  *grpc.Server
	regCollect  []prometheus.Collector
}

var zapLogger *zap.Logger
var customFunc grpc_zap.CodeToLevel

func NewBaseMetrics() *InitMetrics {
	opts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}

	zapLogger, _ = zap.NewDevelopment()
	grpc_zap.ReplaceGrpcLoggerV2(zapLogger)

	return &InitMetrics{
		Reg: reg,
		GrpcServer: grpc.NewServer(
			grpc_middleware.WithUnaryServerChain(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_zap.UnaryServerInterceptor(zapLogger, opts...),
			),
			grpc_middleware.WithStreamServerChain(
				grpc_ctxtags.StreamServerInterceptor(),
				grpc_zap.StreamServerInterceptor(zapLogger, opts...),
			),

			//grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			//	grpc_ctxtags.StreamServerInterceptor(),
			//	grpc_opentracing.StreamServerInterceptor(),
			//	grpc_prometheus.StreamServerInterceptor,
			//	//grpc_zap.StreamServerInterceptor(zapLogger),
			//	//grpc_auth.StreamServerInterceptor(myAuthFunction),
			//	grpc_recovery.StreamServerInterceptor(),
			//)),
			//grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			//	grpc_ctxtags.UnaryServerInterceptor(),
			//	grpc_opentracing.UnaryServerInterceptor(),
			//	grpc_prometheus.UnaryServerInterceptor,
			//	//grpc_zap.UnaryServerInterceptor(zapLogger),
			//	//grpc_auth.UnaryServerInterceptor(myAuthFunction),
			//	grpc_recovery.UnaryServerInterceptor(),
			//)),
		),
	}
}

func (i *InitMetrics) GetGrpcServer() *grpc.Server {
	return i.GrpcServer
}

func (i *InitMetrics) InitAndServe() {
	grpcMetrics.InitializeMetrics(i.GetGrpcServer())

	i.Reg.MustRegister(i.regCollect...)

	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 9092)}

	// Start http server for prometheus.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.GetLogger().Error(err)
		}
	}()
}

func (i *InitMetrics) Registry(cs ...prometheus.Collector) {
	for _, c := range cs {
		i.regCollect = append(i.regCollect, c)
	}
}

func (i *InitMetrics) NewCounterForm(opts prometheus.CounterOpts, labelNames []string) *prometheus.CounterVec {
	vec := prometheus.NewCounterVec(opts, labelNames)
	i.Registry(vec)
	return vec
}

func (i *InitMetrics) NewSummaryForm(opts prometheus.SummaryOpts, labelNames []string) *prometheus.SummaryVec {
	vec := prometheus.NewSummaryVec(opts, labelNames)

	i.Registry(vec)
	return vec
}
