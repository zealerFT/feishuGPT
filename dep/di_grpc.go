package dep

import (
	"feishu/grpc"
)

func DIGRPCServer() (out *grpc.Servlet) {
	container := DI()
	// 启动当前GRPC的服务端
	if err := container.Provide(grpc.NewGrpcServlet); err != nil {
		panic(err)
	}
	if err := container.Invoke(func(s *grpc.Servlet) { out = s }); err != nil {
		panic(err)
	}

	return
}
