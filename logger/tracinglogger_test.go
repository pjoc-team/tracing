package logger

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/uber/jaeger-client-go"
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

func TestSpanContextLog(t *testing.T) {
	if err := tracing.InitOnlyTracingLog("test"); err != nil {
		panic(err.Error())
	}
	tracer, _ := jaeger.NewTracer("test", jaeger.SamplerV2Base{}, jaeger.NewInMemoryReporter())
	opentracing.SetGlobalTracer(tracer)
	span := opentracing.StartSpan("test")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	defer opentracing.SpanFromContext(ctx).Finish()
	parentLog := ContextLog(ctx)
	parentLog.Info("parent")

	ctx, logger := SpanContextLog(ctx)
	logger.Info("child")

	ctx, logger = SpanContextLog(ctx)
	logger.Info("grandson")

	// output:
	// 2021-01-31T17:18:43.279751+08:00 INFO [test,65ac9e65904a0038,65ac9e65904a0038,0,0][44644]parent
	// 2021-01-31T17:18:43.279851+08:00 INFO [test,65ac9e65904a0038,3e44edacfccf7262,65ac9e65904a0038,0][44644]child
	// 2021-01-31T17:18:43.279857+08:00 INFO [test,65ac9e65904a0038,3c453b7c80366692,3e44edacfccf7262,0][44644]grandson
}

func TestContextLog(t *testing.T) {
	if err := tracing.InitOnlyTracingLog("test"); err != nil {
		panic(err.Error())
	}
	tracer, _ := jaeger.NewTracer("test", jaeger.SamplerV2Base{}, jaeger.NewInMemoryReporter())
	opentracing.SetGlobalTracer(tracer)
	span := opentracing.StartSpan("test")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	defer opentracing.SpanFromContext(ctx).Finish()
	parentLog := ContextLog(ctx)

	parentLog.Info("parent") // same
	parentLog.Info("child")  // same

	// output:
	// 2021-01-31T17:18:22.587258+08:00 INFO [test,2f4bbf2967a19f71,2f4bbf2967a19f71,0,0][44447]parent
	// 2021-01-31T17:18:22.58735+08:00 INFO [test,2f4bbf2967a19f71,2f4bbf2967a19f71,0,0][44447]child
}
