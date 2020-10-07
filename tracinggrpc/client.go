package tracinggrpc

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pjoc-team/tracing/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TracingClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		tracer := opentracing.GlobalTracer()
		span, ctx := opentracing.StartSpanFromContext(ctx, method)
		ext.SpanKindRPCClient.Set(span)
		defer span.Finish()
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}
		mdWriter := MDReaderWriter{md}
		err := tracer.Inject(span.Context(), opentracing.TextMap, mdWriter)
		if err != nil {
			logger.ContextLog(ctx).Errorf("inject to metadata err %v", err)
		}
		ctx = metadata.NewOutgoingContext(ctx, md)
		err = invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			ext.Error.Set(span, true)
		}
		return err
	}
}
