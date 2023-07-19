package grpc

import (
	larkBiz "feishu/biz/lark"
	"feishu/config"

	grpcServer "feishu/pkg/grpc"

	larkpb "feishu/proto/go_proto"

	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Dependency struct {
	dig.In

	Cfg  *config.AppConfig
	Lark larkBiz.Lark
}

type Servlet struct {
	server *grpcServer.Servlet
	Dep    *Dependency
	larkpb.UnimplementedLarkServiceServer
}

// Run start run service
func (s *Servlet) Run() {
	s.server.Run()
}

func (s *Servlet) Close() error {
	return s.server.Close()
}

func NewGrpcServlet(cfg *config.AppConfig, dep Dependency) *Servlet {
	maxSize := 20 * 1024 * 1024
	srv := grpcServer.New(&grpcServer.Config{
		ServerEnv:   cfg.AppEnv,
		ServerName:  cfg.AppName,
		ServerHost:  cfg.GRPCServerAddr,
		MetricsHost: cfg.GRPCMetricsAddr,
	}, grpc.MaxRecvMsgSize(maxSize), grpc.MaxSendMsgSize(maxSize))

	servlet := &Servlet{
		server: srv,
		Dep:    &dep,
	}

	larkpb.RegisterLarkServiceServer(srv.Server(), servlet)
	reflection.Register(srv.Server())
	return servlet
}
