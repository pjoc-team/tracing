package util

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func GetTraceID(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		sc, ok := span.Context().(jaeger.SpanContext)
		if ok {
			return sc.TraceID().String()
		}
	}
	return ""
}

func GetTraceInfo(ctx context.Context) (traceID, spanID, parentID string) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		sc, ok := span.Context().(jaeger.SpanContext)
		if ok {
			traceID = sc.TraceID().String()
			spanID = sc.SpanID().String()
			parentID = sc.ParentID().String()
		}
	}
	return traceID, spanID, parentID
}
