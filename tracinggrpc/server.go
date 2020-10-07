package tracinggrpc

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pjoc-team/tracing/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TracingServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		tracer := opentracing.GlobalTracer()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		spanContext, err := tracer.Extract(opentracing.TextMap, MDReaderWriter{md})
		var span opentracing.Span
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			tracinglogger.Log().Errorf("extract from metadata err: %v", err)
		} else {
			span = tracer.StartSpan(
				info.FullMethod,
				ext.RPCServerOption(spanContext),
				opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
				ext.SpanKindRPCServer,
			)
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
		resp, err = handler(ctx, req)
		if err != nil && span != nil {
			ext.Error.Set(span, true)
		}
		return resp, err
	}
}
