package logger

import (
	"fmt"
	"strings"
)

const (
	CtxFieldKeyServiceName  = "ctx.serviceName"
	CtxFieldKeyTraceID      = "ctx.traceID"
	CtxFieldKeySpanID       = "ctx.spanID"
	CtxFieldKeyParentSpanID = "ctx.parentSpanID"
	CtxFieldKeyFlags        = "ctx.flags"
	CtxFieldKeyTime         = "ctx.time"
	CtxFieldKeyLevel        = "ctx.level"
	CtxFieldKeyMsg          = "ctx.msg"
	CtxFieldKeyFile         = "ctx.file"
	CtxFieldKeyRequestID    = "ctx.requestID"
)

type Logger interface {
	TraceID() string

	Debug(args ...interface{})
	Debugln(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infoln(args ...interface{})
	Infof(format string, args ...interface{})

	Print(args ...interface{})
	Println(args ...interface{})
	Printf(format string, args ...interface{})

	Warn(args ...interface{})
	Warnln(args ...interface{})
	Warnf(format string, args ...interface{})

	Warning(args ...interface{})
	Warningln(args ...interface{})
	Warningf(format string, args ...interface{})

	Error(args ...interface{})
	Errorln(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalln(args ...interface{})
	Fatalf(format string, args ...interface{})

	Panic(args ...interface{})
	Panicln(args ...interface{})
	Panicf(format string, args ...interface{})

	IsDebugEnabled() bool

	//support grpclog
	V(l int) bool

	WithField(key string, value interface{}) Logger

	WithFields(fields Fields) Logger
}

type Fields map[string]interface{}

type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid tracinglogger Level: %q", lvl)
}

func (level Level) String() string {
	if b, err := level.MarshalText(); err == nil {
		return string(b)
	} else {
		return "unknown"
	}
}

func (level *Level) UnmarshalText(text []byte) error {
	l, err := ParseLevel(string(text))
	if err != nil {
		return err
	}

	*level = Level(l)

	return nil
}

func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case TraceLevel:
		return []byte("trace"), nil
	case DebugLevel:
		return []byte("debug"), nil
	case InfoLevel:
		return []byte("info"), nil
	case WarnLevel:
		return []byte("warning"), nil
	case ErrorLevel:
		return []byte("error"), nil
	case FatalLevel:
		return []byte("fatal"), nil
	case PanicLevel:
		return []byte("panic"), nil
	}

	return nil, fmt.Errorf("not a valid tracinglogger level %d", level)
}
