package logger

import (
	"context"
	"github.com/pjoc-team/tracing/tracing"
	"testing"
	"time"
)

func TestLogrusText(t *testing.T) {
	tracing.InitOnlyTracingLog("TestLogrusText")
	SetFormatter(FormatText)
	SetLevel(DebugLevel)
	SetReportCallerLevel(InfoLevel)

	ctx := context.Background()
	log, _ := build(ctx, nil)

	log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Debugf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Debugln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Info("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Infof("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Infoln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Print("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Printf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Println("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Warn("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Warnf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Warnln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Warning("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Warningf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Warningln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Error("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Errorf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Errorln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log = log.WithField("key1", "value1")

	log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	fs := make(map[string]interface{})
	fs["key1"] = "value1New"
	fs["key2"] = "value2"
	log = log.WithFields(fs)
	log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	//log.Panic("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	//log.Fatal("failed to fetch URL ", "https://github.com/sirupsen/logrus")
}

func TestReportCallerLevel(t *testing.T) {
	tracing.InitOnlyTracingLog("TestLogrusText")
	SetFormatter(FormatText)
	SetLevel(DebugLevel)
	MinReportCallerLevel(WarnLevel)

	ctx := context.Background()
	log, _ := build(ctx, nil)

	log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Debugf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Debugln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Info("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Infof("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Infoln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Print("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Printf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Println("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Warn("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Warnf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Warnln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Warning("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Warningf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Warningln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Error("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Errorf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Errorln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log = log.WithField("key1", "value1")

	log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	fs := make(map[string]interface{})
	fs["key1"] = "value1New"
	fs["key2"] = "value2"
	log = log.WithFields(fs)
	log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	//log.Panic("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	//log.Fatal("failed to fetch URL ", "https://github.com/sirupsen/logrus")
}

func TestLogrusJson(t *testing.T) {
	tracing.InitOnlyTracingLog("TestLogrusJson")
	SetFormatter(FormatJson)
	SetLevel(DebugLevel)
	MinReportCallerLevel(InfoLevel)

	ctx := context.Background()
	log, _ := build(ctx, nil)

	log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Debugf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Debugln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Info("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Infof("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Infoln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Print("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Printf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Println("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Warn("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Warnf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Warnln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Warning("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Warningf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Warningln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log.Error("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Errorf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	log.Errorln("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	log = log.WithField("key1", "value1")

	log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")

	fs := make(map[string]interface{})
	fs["key1"] = "value1New"
	fs["key2"] = "value2"
	log = log.WithFields(fs)
	log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	//log.Panic("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	//log.Fatal("failed to fetch URL ", "https://github.com/sirupsen/logrus")
}

func TestLogrusJsonCaller(t *testing.T) {
	tracing.InitOnlyTracingLog("TestLogrusJson")
	SetFormatter(FormatJson)
	SetLevel(DebugLevel)
	MinReportCallerLevel(WarnLevel)
	ctx := context.Background()
	log, _ := build(ctx, nil)
	log.Warn("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Info("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	log.Error("failed to fetch URL ", "https://github.com/sirupsen/logrus")
}

/*****************************Benchmark*****************************/

func BenchmarkLogrusText(b *testing.B) {
	tracing.InitOnlyTracingLog("BenchmarkLogrusText")
	SetFormatter(FormatText)
	SetLevel(DebugLevel)
	ctx := context.Background()
	ctx = context.WithValue(ctx, tracing.SpanTagKeyHttpRequestID, time.Now().Format("20060102150405"))
	log, _ := build(ctx, nil)
	for i := 0; i < b.N; i++ {
		log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	}
}

func BenchmarkLogrusJson(b *testing.B) {
	tracing.InitOnlyTracingLog("BenchmarkLogrusJson")
	SetLevel(DebugLevel)
	SetFormatter(FormatJson)
	ctx := context.Background()
	ctx = context.WithValue(ctx, tracing.SpanTagKeyHttpRequestID, time.Now().Format("20060102150405"))
	log, _ := build(ctx, nil)
	for i := 0; i < b.N; i++ {
		log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	}
}

func BenchmarkLogrusTextCaller(b *testing.B) {
	tracing.InitOnlyTracingLog("BenchmarkLogrusTextCaller")
	SetFormatter(FormatText)
	SetLevel(DebugLevel)
	SetReportCallerLevel(DebugLevel)
	ctx := context.Background()
	log, _ := build(ctx, nil)
	for i := 0; i < b.N; i++ {
		log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	}
}

func BenchmarkLogrusJsonCaller(b *testing.B) {
	tracing.InitOnlyTracingLog("BenchmarkLogrusJsonCaller")
	SetLevel(DebugLevel)
	SetFormatter(FormatJson)
	SetReportCallerLevel(DebugLevel)
	ctx := context.Background()
	log, _ := build(ctx, nil)
	for i := 0; i < b.N; i++ {
		log.Debug("failed to fetch URL ", "https://github.com/sirupsen/logrus")
	}
}

/*****************************Benchmark Reflex*****************************/
func BenchmarkLogrusTextReflex(b *testing.B) {
	tracing.InitOnlyTracingLog("BenchmarkLogrusTextReflex")
	SetFormatter(FormatText)
	SetLevel(DebugLevel)
	ctx := context.Background()
	log, _ := build(ctx, nil)
	for i := 0; i < b.N; i++ {
		log.Debugf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	}
}

func BenchmarkLogrusJsonReflex(b *testing.B) {
	tracing.InitOnlyTracingLog("BenchmarkLogrusJsonReflex")
	SetLevel(DebugLevel)
	SetFormatter(FormatJson)
	ctx := context.Background()
	log, _ := build(ctx, nil)
	for i := 0; i < b.N; i++ {
		log.Debugf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	}
}

func BenchmarkLogrusTextCallerReflex(b *testing.B) {
	tracing.InitOnlyTracingLog("BenchmarkLogrusTextCallerReflex")
	SetFormatter(FormatText)
	SetLevel(DebugLevel)
	SetReportCallerLevel(DebugLevel)
	ctx := context.Background()
	log, _ := build(ctx, nil)
	for i := 0; i < b.N; i++ {
		log.Debugf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	}
}

func BenchmarkLogrusJsonCallerReflex(b *testing.B) {
	tracing.InitOnlyTracingLog("BenchmarkLogrusJsonCallerReflex")
	SetLevel(DebugLevel)
	SetFormatter(FormatJson)
	SetReportCallerLevel(DebugLevel)
	ctx := context.Background()
	log, _ := build(ctx, nil)
	for i := 0; i < b.N; i++ {
		log.Debugf("failed to fetch URL %s", "https://github.com/sirupsen/logrus")
	}
}
