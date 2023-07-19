package grpc

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"feishu/pkg/healthcheck"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Config struct {
	ServerEnv       string `json:"server_env"`
	ServerName      string `json:"server_name"`
	HealthCheckAddr string `json:"health_check_addr"`
	// Kubernetes readinessProbe checks GET /health and after (failureThreshold * periodSecond) it stops redirecting traffic to the app (because it continuously returns 500)
	SignalTermWait time.Duration `json:"signal_term_wait"`
	ServerHost     string        `json:"server_host"`
	MetricsHost    string        `json:"metrics_host"`
	SentryDSN      string        `json:"sentry_dsn"`
}

type Servlet struct {
	server     *grpc.Server
	httpHealth *http.Server
	grpcHealth *health.Server
	cfg        *Config
	log        *zap.Logger
}

// Server return grpc.Server for register
func (s *Servlet) Server() *grpc.Server {
	return s.server
}

// New gRPC Servlet
func New(cfg *Config, opts ...grpc.ServerOption) *Servlet {
	lg, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.SignalTermWait == 0 {
		cfg.SignalTermWait = 10 * time.Second
	}
	server := newGRPCServer(lg, opts...)
	servlet := &Servlet{
		cfg:    cfg,
		log:    lg,
		server: server,
	}

	return servlet
}

func (s *Servlet) SetServer(server *grpc.Server) {
	s.server = server
}

func newGRPCServer(l *zap.Logger, opts ...grpc.ServerOption) *grpc.Server {

	zapOptions := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}

	defaultOpts := []grpc.ServerOption{
		grpc.ChainStreamInterceptor(
			grpcctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(l, zapOptions...),
			grpc_prometheus.StreamServerInterceptor,
		),
		grpc.ChainUnaryInterceptor(
			grpcctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(l, zapOptions...),
			grpc_prometheus.UnaryServerInterceptor,
		),
	}

	var serOpts []grpc.ServerOption
	if opts != nil {
		serOpts = append(defaultOpts, opts...)
	} else {
		serOpts = defaultOpts
	}
	return grpc.NewServer(serOpts...)
}

// Run start run service
func (s *Servlet) Run() {
	s.startPrometheus()
	s.startServe()
	s.startHealthMonitor()
	s.waitSignal()
}

// Close servlet
func (s *Servlet) Close() error {
	s.server.GracefulStop()
	return s.log.Sync()
}

func (s *Servlet) startServe() {
	s.log.Info(fmt.Sprintf("gRPC listen %s, running in %s mode", s.cfg.ServerHost, s.cfg.ServerEnv))

	listen, err := net.Listen("tcp", s.cfg.ServerHost)
	if err != nil {
		s.log.Fatal(fmt.Sprintf("failed to listen: %v", err))
	}

	go func() {
		if err := s.server.Serve(listen); err != nil {
			s.log.Fatal(fmt.Sprintf("grpc server listen err: %s\n", err))
		}
	}()
}

func (s *Servlet) startPrometheus() {
	grpc_prometheus.Register(s.server)
	grpc_prometheus.EnableHandlingTimeHistogram()
	http.Handle("/metrics", promhttp.Handler())
	s.log.Info(fmt.Sprintf("gRPC Prometheus listen on %s", s.cfg.MetricsHost))
	go func() {
		if err := http.ListenAndServe(s.cfg.MetricsHost, nil); err != nil {
			s.log.Fatal(fmt.Sprintf("Failed to start prometheus server: %s", err))
		}
	}()
}

func (s *Servlet) startHealthMonitor() {
	hs := health.NewServer()
	hs.SetServingStatus(s.cfg.ServerName, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s.server, hs)
	s.grpcHealth = hs

	server := &http.Server{Addr: s.cfg.HealthCheckAddr, Handler: healthcheck.NewHandler()}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return
			}
			s.log.Error("start health server failed", zap.Error(err))
		}
	}()
	s.httpHealth = server

	s.log.Info("gRPC Health boot")
}

func (s *Servlet) waitSignal() {
	quit := make(chan os.Signal, 1)
	signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	s.log.Info("GRPC server receive TERM signal")
	if err := s.httpHealth.Close(); err != nil {
		s.log.Error("health server shutdown", zap.Error(err))
	}
	s.grpcHealth.SetServingStatus(s.cfg.ServerName, healthpb.HealthCheckResponse_NOT_SERVING)
	time.Sleep(s.cfg.SignalTermWait)
}
