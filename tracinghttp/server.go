package tracinghttp

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracinglogger "github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/util"
	"net/http"
)

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		newCtx := r.Context()
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			tracinglogger.Log().Errorf("extract from header err: %v", err)
		} else {
			span := opentracing.GlobalTracer().StartSpan(pattern, ext.RPCServerOption(spanCtx))
			newCtx = opentracing.ContextWithSpan(newCtx, span)
			requestID := r.Header.Get(string(string(tracing.HttpHeaderKeyXRequestID)))
			if requestID != "" {
				span.SetTag(string(tracing.SpanTagKeyHttpRequestID), requestID)
				newCtx = context.WithValue(newCtx, tracing.SpanTagKeyHttpRequestID, requestID)
				w.Header().Add(string(tracing.HttpHeaderKeyXRequestID), requestID)
			}
			defer span.Finish()
		}
		w.Header().Add(string(tracing.TraceID), util.GetTraceID(newCtx))
		handler(w, r.WithContext(newCtx))
	})
}

func TracingServerInterceptor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newCtx := r.Context()
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			tracinglogger.Log().Errorf("extract from header err: %v", err)
		} else {
			span := opentracing.GlobalTracer().StartSpan(r.RequestURI, ext.RPCServerOption(spanCtx))
			newCtx = opentracing.ContextWithSpan(newCtx, span)
			requestID := r.Header.Get(string(tracing.HttpHeaderKeyXRequestID))
			if requestID != "" {
				span.SetTag(string(tracing.SpanTagKeyHttpRequestID), requestID)
				newCtx = context.WithValue(newCtx, tracing.SpanTagKeyHttpRequestID, requestID)
			}
			defer span.Finish()
		}
		h.ServeHTTP(w, r.WithContext(newCtx))
	})
}
