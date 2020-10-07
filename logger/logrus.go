package tracinglogger

import (
	"bytes"
	"errors"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/util"
	"github.com/sirupsen/logrus"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	paramCtxFormat = "ctx.format"
)

var (
	logrusLog        = logrus.New()
	logrusCallerLog  *logrusReportCaller
	logrusFormatType = FormatText
	callerLevels     []logrus.Level
)

func init() {
	logrusLog.SetFormatter(&defaultLoggerFormatter{})
	logrusLog.SetLevel(logrus.GetLevel())
}

type logrusFactory struct{}

func newLogrusReportCaller(level []logrus.Level) *logrusReportCaller {
	callerLevels = level
	logrusCallerLog = &logrusReportCaller{
		level:  level,
		caller: true,
	}
	logrusLog.AddHook(&JaegerCallerHook{})
	return logrusCallerLog
}

type logrusReportCaller struct {
	level  []logrus.Level
	caller bool
}

func (f *logrusFactory) setFormatter(ft FormatType) error {
	logrusFormatType = ft
	switch ft {
	case FormatText:
		logrusLog.SetFormatter(&defaultLoggerFormatter{})
		return nil
	case FormatJson:
		logrusLog.SetFormatter(&defaultLoggerJsonFormatter{})
		return nil
	default:
		logrusLog.SetFormatter(&defaultLoggerFormatter{})
		return errors.New("UNKNOWN FormatType")
	}
}

func (f *logrusFactory) setOutput(output io.Writer) {
	logrusLog.SetOutput(output)
}

func (f *logrusFactory) setLevel(lvl Level) error {
	l, err := logrus.ParseLevel(lvl.String())
	if err == nil {
		logrusLog.SetLevel(l)
	}
	return err
}

func (f *logrusFactory) setReportCallerLevel(lvl ...Level) error {
	if logrusCallerLog == nil {
		levels := make([]logrus.Level, 0, len(logrus.AllLevels))
		for _, v := range lvl {
			l, err := logrus.ParseLevel(v.String())
			if err == nil {
				levels = append(levels, l)
			} else {
				return err
			}
		}
		if len(lvl) == 0 {
			levels = append(levels, logrus.PanicLevel)
		}
		newLogrusReportCaller(levels)
	}
	return nil
}

func (f *logrusFactory) buildLogger(ti traceInfo) Logger {
	if logrusFormatType == FormatJson {
		fs := logrus.Fields{}
		serviceName := tracing.GetServiceName()
		if serviceName != "" {
			fs[CtxFieldKeyServiceName] = tracing.GetServiceName()
		}
		if ti.traceID != "" {
			fs[CtxFieldKeyTraceID] = ti.traceID
			fs[CtxFieldKeySpanID] = ti.spanID
			fs[CtxFieldKeyParentSpanID] = ti.parentSpanID
			fs[CtxFieldKeyFlags] = ti.flags
		}
		if ti.requestID != "" {
			fs[CtxFieldKeyRequestID] = ti.requestID
		}
		return newEntry(logrusLog.WithFields(fs))
	}
	fs := logrus.Fields{
		paramCtxFormat:     ti.String(),
		CtxFieldKeyTraceID: ti.traceID,
	}
	return newEntry(logrusLog.WithFields(fs))
}

type defaultLoggerFormatter struct{}

func (formatter defaultLoggerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(time.RFC3339Nano)
	msg := util.ToStr(entry.Data[paramCtxFormat])
	var buf strings.Builder
	buf.WriteString(msg)

	if entry.Caller != nil {
		buf.WriteByte(' ')
		buf.WriteString(util.GetPreAndFileName(entry.Caller.File))
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(entry.Caller.Line))
		buf.WriteByte(' ')
	}

	if len(entry.Data) > 0 {
		buf.WriteString(util.MapToStr(entry.Data, paramCtxFormat, CtxFieldKeyTraceID))
	}
	buf.WriteString(entry.Message)

	var logBuf bytes.Buffer
	logBuf.WriteString(timestamp)
	logBuf.WriteByte(' ')
	logBuf.WriteString(strings.ToUpper(entry.Level.String()))
	logBuf.WriteByte(' ')
	logBuf.WriteString(buf.String())
	logBuf.WriteByte('\n')
	return logBuf.Bytes(), nil
}

type defaultLoggerJsonFormatter struct{}

func (formatter defaultLoggerJsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(map[string]interface{})
	data[CtxFieldKeyTime] = entry.Time.Format(time.RFC3339Nano)
	data[CtxFieldKeyMsg] = entry.Message
	data[CtxFieldKeyLevel] = entry.Level.String()
	for k, v := range entry.Data {
		if k == CtxFieldKeyFile && callerLevels != nil {
			for _, lvl := range callerLevels {
				if lvl == entry.Level {
					data[k] = v
					break
				}
			}
		} else {
			data[k] = v
		}
	}
	return util.GetJsonBytes(data), nil
}

type JaegerCallerHook struct {
}

func (df *JaegerCallerHook) Levels() []logrus.Level {
	return logrusCallerLog.level
}

func (df *JaegerCallerHook) Fire(entry *logrus.Entry) error {
	if logrusCallerLog.caller {
		entry.Caller = util.GetJaegerCaller()
		if entry.Caller != nil && logrusFormatType == FormatJson {
			var buf strings.Builder
			buf.WriteString(util.GetPreAndFileName(entry.Caller.File))
			buf.WriteByte(':')
			buf.WriteString(strconv.Itoa(entry.Caller.Line))
			entry.Data[CtxFieldKeyFile] = buf.String()
		}
	}
	return nil
}

type entry struct {
	log *logrus.Entry
}

func (e *entry) V(l int) bool {
	// Returns values of debug level.
	return e.IsDebugEnabled()
}

func newEntry(l *logrus.Entry) *entry {
	return &entry{log: l}
}

func (e *entry) TraceID() string {
	return util.ToStr(e.log.Data[CtxFieldKeyTraceID])
}

func (e *entry) Debug(args ...interface{}) {
	e.log.Debug(args...)
}

func (e *entry) Debugln(args ...interface{}) {
	e.log.Debugln(args...)
}

func (e *entry) Debugf(format string, args ...interface{}) {
	e.log.Debugf(format, args...)
}

func (e *entry) Print(args ...interface{}) {
	e.log.Print(args...)
}

func (e *entry) Println(args ...interface{}) {
	e.log.Println(args...)
}

func (e *entry) Printf(format string, args ...interface{}) {
	e.log.Printf(format, args...)
}

func (e *entry) Info(args ...interface{}) {
	e.log.Info(args...)
}

func (e *entry) Infoln(args ...interface{}) {
	e.log.Infoln(args...)
}

func (e *entry) Infof(format string, args ...interface{}) {
	e.log.Infof(format, args...)
}

func (e *entry) Warn(args ...interface{}) {
	e.log.Warn(args...)
}

func (e *entry) Warnln(args ...interface{}) {
	e.log.Warnln(args...)
}

func (e *entry) Warnf(format string, args ...interface{}) {
	e.log.Warnf(format, args...)
}

func (e *entry) Warning(args ...interface{}) {
	e.log.Warning(args...)
}

func (e *entry) Warningln(args ...interface{}) {
	e.log.Warningln(args...)
}

func (e *entry) Warningf(format string, args ...interface{}) {
	e.log.Warningf(format, args...)
}

func (e *entry) Error(args ...interface{}) {
	e.log.Error(args...)
}

func (e *entry) Errorln(args ...interface{}) {
	e.log.Errorln(args...)
}

func (e *entry) Errorf(format string, args ...interface{}) {
	e.log.Errorf(format, args...)
}

func (e *entry) Fatal(args ...interface{}) {
	e.log.Fatal(args...)
}

func (e *entry) Fatalln(args ...interface{}) {
	e.log.Fatalln(args...)
}

func (e *entry) Fatalf(format string, args ...interface{}) {
	e.log.Fatalf(format, args...)
}

func (e *entry) Panic(args ...interface{}) {
	e.log.Panic(args...)
}

func (e *entry) Panicln(args ...interface{}) {
	e.log.Panicln(args...)
}

func (e *entry) Panicf(format string, args ...interface{}) {
	e.log.Panicf(format, args...)
}

func (e *entry) IsDebugEnabled() bool {
	return logrusLog.Level >= logrus.DebugLevel
}

func (e *entry) String() (string, error) {
	return e.log.String()
}

func (e *entry) WithField(key string, value interface{}) Logger {
	return &entry{log: e.log.WithField(key, value)}
}

func (e *entry) WithFields(fields Fields) Logger {
	return &entry{log: e.log.WithFields(logrus.Fields(fields))}
}
