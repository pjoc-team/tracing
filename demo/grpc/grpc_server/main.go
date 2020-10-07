// Package main implements a server for Greeter service.
package main

import (
	"context"
	"github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/tracinggrpc"
	"github.com/pjoc-team/tracing/tracinghttp"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log := tracinglogger.ContextLog(ctx)
	log.Printf("Received: %v", in.GetName())
	_, res, err := tracinghttp.Get(ctx, http.DefaultClient, "http://localhost:8084/sayHello")
	if err != nil {
		log.Errorf("%v", err)
	} else {
		result, _ := ioutil.ReadAll(res.Body)
		log.Infof("res:%s", result)
		res.Body.Close()
	}
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	tracing.InitOnlyTracingLog("grpc_server")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(tracinggrpc.TracingServerInterceptor()))
	s := grpc.NewServer(opts...)
	grpclog.SetLoggerV2(tracinglogger.Log())
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
