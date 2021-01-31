package tracinghttp

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracinglogger "github.com/pjoc-team/tracing/logger"
	"net/http"
)

type tracingClient struct {
	next http.RoundTripper
}

func (h *tracingClient) RoundTrip(request *http.Request) (*http.Response, error) {
	uri := request.URL.Path
	ctx := request.Context()
	span, newCtx := opentracing.StartSpanFromContext(ctx, uri)
	defer span.Finish()

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, request.URL.String())
	ext.HTTPMethod.Set(span, request.Method)
	err := span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header),
	)
	log := tracinglogger.ContextLog(newCtx)
	if err != nil {
		log.Errorf("inject to HTTPHeaders err %v", err)
	}
	if log.IsDebugEnabled() {
		log.Debugf("header: %v", request.Header)
	}

	req := request.WithContext(newCtx)

	resp, err := h.next.RoundTrip(req)
	if err != nil {
		ext.Error.Set(span, true)
	} else {
		ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))
	}
	return resp, err
}

// NewClient 创建支持tracing的http client
func NewClient(r *http.Client) *http.Client {
	if r == nil {
		r = http.DefaultClient
	}
	if r.Transport == nil {
		r.Transport = http.DefaultTransport
	}

	// metrics
	h := &tracingClient{next: r.Transport}

	// replace transport and copy other fields from r
	c := &http.Client{
		Transport:     h,
		CheckRedirect: r.CheckRedirect,
		Jar:           r.Jar,
		Timeout:       r.Timeout,
	}

	return c
}
