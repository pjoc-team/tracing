package tracing

// Key context key
type Key string

const (
	HttpHeaderKeyXRequestID Key = "X-Request-Id"
	SpanTagKeyHttpRequestID Key = "http.request_id"
	TraceID                 Key = "Trace-Id" // TraceID http响应header内返回的traceID
)
