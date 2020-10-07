package tracing

const (
	HttpHeaderKeyXRequestID = "X-Request-Id"
	SpanTagKeyHttpRequestID = "http.request_id"
	TraceID                 = "Trace-Id" // TraceID http响应header内返回的traceID
)
