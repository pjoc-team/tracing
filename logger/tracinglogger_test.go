package tracinglogger

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pjoc-team/tracing/tracing"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	SetLevel(WarnLevel)
	log := ContextLog(context.Background())
	log = log.WithField("k1", "value")
	log.Warn("test log")
	log.Info("test log")
}

func TestReportCallerLog(t *testing.T) {
	SetLevel(WarnLevel)
	SetReportCallerLevel(WarnLevel)
	log := ContextLog(context.Background())
	log = log.WithField("k1", "value")
	log.Warn("test log")
	log.Info("test log")
}

func build(ctx context.Context, s opentracing.Span) (Logger, opentracing.Span) {
	tracer := opentracing.GlobalTracer()
	var span opentracing.Span
	if s != nil {
		span = tracer.StartSpan(
			"TestTracerLogger",
			opentracing.Tag{Key: string(ext.Component), Value: "test"},
			ext.SpanKindRPCClient,
			opentracing.ChildOf(s.Context()),
		)
	} else {
		span = tracer.StartSpan(
			"TestTracerLogger",
			opentracing.Tag{Key: string(ext.Component), Value: "test"},
			ext.SpanKindRPCClient,
		)
	}
	ctx = opentracing.ContextWithSpan(ctx, span)
	return ContextLog(ctx), span
}

func TestSpanIncludeSpan(t *testing.T) {
	os.Setenv("JAEGER_SERVICE_NAME", "JaegerTest")
	tracing.InitOnlyTracingLog("")
	ctx := context.Background()
	t1, span := build(ctx, nil)
	t1.Debug("t1-->test")
	t1.Print("t1-->test")
	t1.Info("t1-->test")
	t1.Warn("t1-->Out of memory")
	t1.Warning("t1-->Out of memory")
	t1.Error("t1-->NullPointException")
	t1.Printf("t1-->test %d", 1)
	t1 = t1.WithField("k1", "value1")
	t1 = t1.WithField("k2", "value2")
	f := Fields{}
	f["test1"] = "test111"
	t1 = t1.WithFields(f)
	t1.Info("test log")
	t2, _ := build(ctx, span)
	t2.Debugln("t2-->pjoc 666")
	t2.Println("t2-->pjoc 666")
	t2.Infoln("t2-->pjoc 666")
	t2.Warnln("t2-->Out of memory")
}

func TestConfig(t *testing.T) {
	c := &tracing.Config{
		ServiceName:      "JaegerTest",
		SamplerRateLimit: 1,
		JaegerAgent:      "localhost:6381",
	}
	err := tracing.InitFromConfig(c)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	ctx := context.Background()
	t1, _ := build(ctx, nil)
	t1.Debug("t1-->test")
	t1.Print("t1-->test")
	t1.Info("t1-->test")
	t1.Warn("t1-->Out of memory")
	t1.Warning("t1-->Out of memory")
}

func TestJson(t *testing.T) {
	os.Setenv("JAEGER_SERVICE_NAME", "JaegerTest")
	tracing.InitOnlyTracingLog("")
	SetFormatter(FormatJson)
	ctx := context.Background()
	t1, span := build(ctx, nil)
	t1.Debug("t1-->test")
	t1.Print("t1-->test")
	t1.Info("t1-->test")
	t1.Warn("t1-->Out of memory")
	t1.Warning("t1-->Out of memory")
	t1.Error("t1-->NullPointException")
	t1.Printf("t1-->test %d", 1)
	t1 = t1.WithField("k1", "value1")
	t1 = t1.WithField("k2", "value2")
	f := Fields{}
	f["test1"] = "test111"
	t1 = t1.WithFields(f)
	t1.Info("test log")
	t2, _ := build(ctx, span)
	t2.Debugln("t2-->pjoc 666")
	t2.Println("t2-->pjoc 666")
	t2.Infoln("t2-->pjoc 666")
	t2.Warnln("t2-->Out of memory")
}

func TestReportCallerJson(t *testing.T) {
	os.Setenv("JAEGER_SERVICE_NAME", "JaegerTest")
	tracing.InitOnlyTracingLog("")
	SetFormatter(FormatJson)
	SetLevel(DebugLevel)
	SetReportCallerLevel(AllLevels...)
	ctx := context.Background()
	t1, span := build(ctx, nil)
	t1.Debug("t1-->test")
	t1.Print("t1-->test")
	t1.Info("t1-->test")
	t1.Warn("t1-->Out of memory")
	t1.Warning("t1-->Out of memory")
	t1.Error("t1-->NullPointException")
	t1.Printf("t1-->test %d", 1)
	t1 = t1.WithField("k1", "value1")
	t1 = t1.WithField("k2", "value2")
	f := Fields{}
	f["test1"] = "test111"
	t1 = t1.WithFields(f)
	t1.Info("test log")
	t2, _ := build(ctx, span)
	t2.Debugln("t2-->pjoc 666")
	t2.Println("t2-->pjoc 666")
	t2.Infoln("t2-->pjoc 666")
	t2.Warnln("t2-->Out of memory")
}

func TestSetReportCaller(t *testing.T) {
	os.Setenv("JAEGER_SERVICE_NAME", "JaegerTest")
	tracing.InitOnlyTracingLog("")
	SetLevel(DebugLevel)
	SetReportCallerLevel(DebugLevel)
	ctx := context.Background()
	t1, span := build(ctx, nil)
	t1.Debug("t1-->test")
	t1.Print("t1-->test")
	t1.Info("t1-->test")
	t1.Warn("t1-->Out of memory")
	t1.Warning("t1-->Out of memory")
	t1.Error("t1-->NullPointException")
	t1.Printf("t1-->test %d", 1)
	t1 = t1.WithField("k1", "value1")
	t1 = t1.WithField("k2", "value2")
	f := Fields{}
	f["test1"] = "test111"
	t1 = t1.WithFields(f)
	t1.Info("test log")
	t2, _ := build(ctx, span)
	t2.Debugln("t2-->pjoc 666")
	t2.Println("t2-->pjoc 666")
	t2.Infoln("t2-->pjoc 666")
	t2.Warnln("t2-->Out of memory")
}

func TestNoTracingLog(t *testing.T) {
	SetLevel(DebugLevel)
	SetReportCallerLevel(DebugLevel)
	Log().Debug("hello madam")
	Log().Info("hello madam")
}

func TestNoTracingLogJson(t *testing.T) {
	SetLevel(DebugLevel)
	SetReportCallerLevel(DebugLevel)
	SetFormatter(FormatJson)
	Log().Debug("hello madam")
	Log().Info("hello madam")
}

func TestIsDebugEnabled(t *testing.T) {
	os.Setenv("JAEGER_SERVICE_NAME", "JaegerTest")
	SetLevel(TraceLevel)
	SetReportCallerLevel(InfoLevel)
	Log().Info("hello madam")
	if Log().IsDebugEnabled() {
		Log().Info("IsDebugEnabled true")
		Log().Debug("hello madam")
	}
}

func TestEntry_TraceID(t *testing.T) {
	Log().Info(Log().TraceID())

	os.Setenv("JAEGER_SERVICE_NAME", "JaegerTest")
	tracing.InitOnlyTracingLog("")
	SetLevel(DebugLevel)
	SetReportCallerLevel(DebugLevel)
	ctx := context.Background()
	t1, _ := build(ctx, nil)

	t1.Info(t1.TraceID())

	SetFormatter(FormatJson)

	t1, _ = build(ctx, nil)

	t1.Info(t1.TraceID())
}
