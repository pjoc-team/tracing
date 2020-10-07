// Package main implements a client for Greeter service.
package main

import (
	"context"
	"github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/tracinggrpc"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	tracing.InitOnlyTracingLog("grpc_client")
	tracing.HandleFunc(func(ctx context.Context) {
		log := tracinglogger.ContextLog(ctx)
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())
		opts = append(opts, grpc.WithBlock())
		//初始化grpc client拦截器
		opts = append(opts, grpc.WithUnaryInterceptor(tracinggrpc.TracingClientInterceptor()))
		conn, err := grpc.Dial(address, opts...)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewGreeterClient(conn)
		name := defaultName
		if len(os.Args) > 1 {
			name = os.Args[1]
		}
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	})
}
