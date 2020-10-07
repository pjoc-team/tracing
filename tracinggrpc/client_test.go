package tracinggrpc

import (
	"google.golang.org/grpc"
)

func ExampleTracingClientInterceptor() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	//初始化grpc client拦截器
	opts = append(opts, grpc.WithUnaryInterceptor(TracingClientInterceptor()))
	grpc.Dial("localhost:50051", opts...)
}
