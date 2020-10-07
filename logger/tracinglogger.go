package tracinglogger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/util"
	"github.com/uber/jaeger-client-go"
	"io"
	"sync"
)

var (
	log  Logger
	once sync.Once
)

func SetOutput(output io.Writer) {
	getFactory().setOutput(output)
}

//设置日志格式化方式
func SetFormatter(ft FormatType) error {
	return getFactory().setFormatter(ft)
}

//设置日志级别
func SetLevel(lvl Level) error {
	return getFactory().setLevel(lvl)
}

//设置需要打印文件名行号的日志级别
func SetReportCallerLevel(lvl ...Level) error {
	return getFactory().setReportCallerLevel(lvl...)
}

//设置需要打印文件名行号的最小日志级别
func MinReportCallerLevel(lvl Level) error {
	arr := make([]Level, 1)
	for _, v := range AllLevels {
		if v <= lvl {
			arr = append(arr, v)
		}
	}
	return getFactory().setReportCallerLevel(arr...)
}

//获取带有traceing信息的Logger对象
func ContextLog(ctx context.Context) Logger {
	return getFactory().buildLogger(buildTraceInfo(ctx))
}

//非tracing环境,Logger对象,简化ContextLog(context.Background())操作
func Log() Logger {
	if log == nil {
		once.Do(func() {
			log = ContextLog(context.Background())
		})
	}
	return log
}

//构建trace信息
func buildTraceInfo(ctx context.Context) traceInfo {
	//traceID, spanID, parentSpanID, flags, requestID := "", "", "", "0", ""
	ti := traceInfo{flags: "0"}
	if ctx != nil {
		span := opentracing.SpanFromContext(ctx)
		if span != nil {
			sc, ok := span.Context().(jaeger.SpanContext)
			if ok {
				ti.traceID = sc.TraceID().String()
				ti.spanID = sc.SpanID().String()
				ti.parentSpanID = sc.ParentID().String()
				if sc.IsSampled() {
					ti.flags = "1"
				}
			}
		}
		ti.requestID = util.ToStr(ctx.Value(tracing.SpanTagKeyHttpRequestID))
	}
	return ti
}
