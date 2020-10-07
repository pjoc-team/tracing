package logger

import (
	"github.com/pjoc-team/tracing/tracing"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	factory *Factory
)

func init() {
	factory = &Factory{Handler: &logrusFactory{}}
}

type Factory struct {
	Handler
}

type traceInfo struct {
	traceID, spanID, parentSpanID, flags, requestID string
}

func (ti *traceInfo) String() string {
	var buf strings.Builder
	if ti.requestID != "" {
		buf.WriteByte('[')
		buf.WriteString(ti.requestID)
		buf.WriteByte(']')
	}
	buf.WriteByte('[')
	buf.WriteString(tracing.GetServiceName())
	if ti.traceID != "" {
		buf.WriteByte(',')
		buf.WriteString(ti.traceID)
		buf.WriteByte(',')
		buf.WriteString(ti.spanID)
		buf.WriteByte(',')
		buf.WriteString(ti.parentSpanID)
		buf.WriteByte(',')
		buf.WriteString(ti.flags)
	}
	buf.WriteByte(']')
	if buf.String() == "[]" {
		buf.Reset()
	}
	buf.WriteByte('[')
	buf.WriteString(strconv.Itoa(os.Getpid()))
	buf.WriteByte(']')
	return buf.String()
}

type Handler interface {
	setOutput(output io.Writer)
	setFormatter(ft FormatType) error
	setLevel(lvl Level) error
	setReportCallerLevel(lvl ...Level) error
	buildLogger(ti traceInfo) Logger
}

func getFactory() *Factory {
	return factory
}
