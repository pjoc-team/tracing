package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

//使普通方法在trace上下文中，一般用在第一个发起者
func HandleFunc(handler func(ctx context.Context)) {
	span := opentracing.GlobalTracer().StartSpan("HandleFunc")
	handler(opentracing.ContextWithSpan(context.Background(), span))
	defer span.Finish()
}

//手动上报
func Forward(ctx context.Context, operationName string, handler func(ctx context.Context)) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	handler(opentracing.ContextWithSpan(ctx, span))
	defer span.Finish()
}

//复杂场景手动上报
func Start(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctxNew := opentracing.StartSpanFromContext(ctx, operationName)
	return span, opentracing.ContextWithSpan(ctxNew, span)
}

func BuildContextByCarrier(carrier opentracing.TextMapCarrier, operationName string, tag ...string) context.Context {
	tracer := opentracing.GlobalTracer()
	extractedContext, err := tracer.Extract(
		opentracing.TextMap,
		carrier,
	)
	if err != nil && err != opentracing.ErrSpanContextNotFound {
		fmt.Println("BuildContextByCarrier extract error:", err.Error())
		return context.Background()
	} else {
		span := tracer.StartSpan(
			operationName,
			opentracing.FollowsFrom(extractedContext),
		)
		if len(tag) > 0 {
			for _, value := range tag {
				ext.MessageBusDestination.Set(span, value)
			}
		}
		defer span.Finish()
		return opentracing.ContextWithSpan(context.Background(), span)
	}
}
