package tracinggin

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracinglogger "github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/util"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func TracingServerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := opentracing.GlobalTracer()
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		ctx, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
		var span opentracing.Span
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			tracinglogger.Log().Errorf("extract from header err: %v", err)
		} else {
			span = tracer.StartSpan(c.Request.Method, ext.RPCServerOption(ctx))
			ext.HTTPMethod.Set(span, c.Request.Method)
			ext.HTTPUrl.Set(span, c.Request.URL.Path)
			ext.Component.Set(span, "net/http")
			requestID := c.Request.Header.Get(string(tracing.HttpHeaderKeyXRequestID))
			newCtx := c.Request.Context()
			if requestID != "" {
				span.SetTag(string(tracing.SpanTagKeyHttpRequestID), requestID)
				newCtx = context.WithValue(newCtx, tracing.SpanTagKeyHttpRequestID, requestID)
			}
			c.Request = c.Request.WithContext(
				opentracing.ContextWithSpan(newCtx, span),
			)
		}
		c.Next()
		if span != nil {
			if c.Writer.Status() >= http.StatusBadRequest {
				ext.Error.Set(span, true)
			}
			ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
			span.Finish()
		}
	}
}

type (
	Option interface {
		apply(*logOptions)
	}

	logOptions struct {
		skipStack    int
		stackRow     int
		skipPaths    []string
		includePaths []string
		respPrefix   string
		panicPrefix  string
	}

	logOption struct {
		o func(*logOptions)
	}
)

var defaultOptions = logOptions{
	skipStack:   3,
	stackRow:    3,
	respPrefix:  "gin-resp-status",
	panicPrefix: "gin-panic-recovered",
}

func (o *logOption) apply(do *logOptions) {
	o.o(do)
}

// PanicStack 设置请求发生异常时，打印需要跳过的堆栈、需要打印的堆栈行数
func PanicStack(skipStack, stackRow int) Option {
	return &logOption{o: func(o *logOptions) {
		o.skipStack = skipStack
		o.stackRow = stackRow
	}}
}

// IncludePath 设置需要打印响应结果的请求路径 默认全包含     * 表任意，可以每级自由定义
func IncludePath(includePaths ...string) Option {
	return &logOption{o: func(o *logOptions) {
		o.includePaths = includePaths
	}}
}

// SkipPath 设置需要跳过打印响应结果的请求路径 默认都不跳过   * 表任意，可以每级自由定义
func SkipPath(skipPaths ...string) Option {
	return &logOption{o: func(o *logOptions) {
		o.skipPaths = skipPaths
	}}
}

// LogPrefix 设置tracing响应结果日志输出的日志前缀关键字
func LogPrefix(resp, panic string) Option {
	return &logOption{o: func(o *logOptions) {
		o.respPrefix = resp
		o.panicPrefix = panic
	}}
}

func checkPath(paths []string, path string, skip bool) bool {
	if paths == nil || len(paths) == 0 {
		if skip {
			return false
		}
		return true
	}
	for _, v := range paths {
		if v == "/*" || v == "*" {
			return true
		}
		ps := strings.Split(path, "/")
		pl := len(ps)
		is := strings.Split(v, "/")
		il := len(is)
		for k, s := range is {
			if k <= (pl - 1) {
				if s == "*" {
					if k == (il - 1) {
						return true
					}
					continue
				}
				if s != ps[k] {
					return false
				}
			} else if s == "*" || s == "" {
				if k == (il - 1) {
					return true
				}
				continue
			}
		}
	}
	return false
}

func path(opt logOptions, url string) bool {
	if !checkPath(opt.includePaths, url, false) {
		return false
	}
	if checkPath(opt.skipPaths, url, true) {
		return false
	}
	return true
}

// TracingLogInterceptor 每次请求结束打印tracing响应结果、耗时、请求url，发生异常时打印堆栈，放在TracingServerInterceptor后调用
func TracingLogInterceptor(opt ...Option) gin.HandlerFunc {
	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				tracinglogger.ContextLog(c.Request.Context()).Errorf("%s|%s|%s|%s|%s|%s",
					opts.panicPrefix,
					c.ClientIP(),
					c.Request.Method,
					c.Request.URL.Path,
					stack(opts.skipStack, opts.stackRow),
					err,
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		if !path(opts, c.Request.URL.Path) {
			return
		}
		// Start timer
		start := time.Now()
		// Process request
		c.Next()
		// Try not to use reflections
		log := &strings.Builder{}
		log.WriteString(opts.respPrefix)
		log.WriteString("|")
		log.WriteString(strconv.Itoa(c.Writer.Status()))
		log.WriteString("|")
		log.WriteString(time.Since(start).String())
		log.WriteString("|")
		log.WriteString(c.ClientIP())
		log.WriteString("|")
		log.WriteString(c.Request.Method)
		log.WriteString("|")
		log.WriteString(c.Request.URL.Path)
		if c.Request.URL.RawQuery != "" {
			log.WriteString("?")
			log.WriteString(c.Request.URL.RawQuery)
		}
		tracinglogger.ContextLog(c.Request.Context()).Info(log.String())
	}
}

//尽可能提高日志可读性，同时不能包含换行符，防止log无法收集
func stack(skipStack, stackRow int) []byte {
	buf := new(bytes.Buffer)
	for i := skipStack; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok || (i-skipStack > (stackRow - 1)) {
			break
		}
		format := "%s:%d"
		if i != skipStack {
			format = ",%s:%d"
		}
		fmt.Fprintf(buf, format, util.GetPreAndFileName(file), line)
	}
	return buf.Bytes()
}
