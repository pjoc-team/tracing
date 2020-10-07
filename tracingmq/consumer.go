package tracingmq

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pjoc-team/tracing/logger"
)

// Deprecated: 推荐使用 tracing.BuildContextByCarrier
func TracingMqConsumer(ctx context.Context, tarceMqData TraceMqData) context.Context {
	tracer := opentracing.GlobalTracer()
	extractedContext, err := tracer.Extract(
		opentracing.TextMap,
		tarceMqData.Carriers,
	)
	if err != nil && err != opentracing.ErrSpanContextNotFound {
		tracinglogger.Log().Errorf("TracingMqConsumer extract error: %s ", err.Error())
		return ctx
	} else {
		span := tracer.StartSpan(
			"mq_consumer",
			opentracing.FollowsFrom(extractedContext),
		)
		ext.MessageBusDestination.Set(span, tarceMqData.Topic)
		defer span.Finish()
		return opentracing.ContextWithSpan(ctx, span)
	}
}
