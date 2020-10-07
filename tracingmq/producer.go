package tracingmq

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pjoc-team/tracing/logger"
)

func TracingMqProducer(ctx context.Context, topic string, data []byte) *TraceMqData {
	span, nctx := opentracing.StartSpanFromContext(ctx, "mq_producer")
	ext.SpanKindProducer.Set(span)
	ext.MessageBusDestination.Set(span, topic)
	defer span.Finish()
	carriers := opentracing.TextMapCarrier{}
	err := span.Tracer().Inject(
		span.Context(),
		opentracing.TextMap,
		carriers)
	if err != nil {
		tracinglogger.ContextLog(nctx).Errorf("TracingMqProducer inject to TextMap err %v", err)
	}
	return &TraceMqData{
		Carriers: carriers,
		Topic:    topic,
		Data:     data,
	}
}
