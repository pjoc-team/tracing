package tracinggrpc

import (
	"github.com/pjoc-team/tracing/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

func ExampleTracingServerInterceptor() {
	net.Listen("tcp", ":50051")
	var opts []grpc.ServerOption
	//初始化grpc服务端拦截器
	opts = append(opts, grpc.UnaryInterceptor(TracingServerInterceptor()))
	grpc.NewServer(opts...)
	//还可以设置grpc log
	grpclog.SetLoggerV2(tracinglogger.Log())
}
