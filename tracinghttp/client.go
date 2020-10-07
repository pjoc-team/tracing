package tracinghttp

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracinglogger "github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/util"
	"io"
	"net/http"
)

func GetDo(ctx context.Context, c *http.Client, url string, header map[string]string) (sctx context.Context, resp *http.Response, err error) {
	return do(ctx, c, url, "GET", "", header, nil)
}

func PostDo(ctx context.Context, c *http.Client, url string, header map[string]string, body io.Reader) (sctx context.Context, resp *http.Response, err error) {
	return do(ctx, c, url, "POST", "", header, body)
}

func Get(ctx context.Context, c *http.Client, url string) (sctx context.Context, resp *http.Response, err error) {
	return do(ctx, c, url, "GET", "", nil, nil)
}

func Post(ctx context.Context, c *http.Client, url, contentType string, body io.Reader) (sctx context.Context, resp *http.Response, err error) {
	return do(ctx, c, url, "POST", contentType, nil, body)
}

func do(ctx context.Context, c *http.Client, url, method, contentType string, header map[string]string, body io.Reader) (spanCtx context.Context, resp *http.Response, err error) {
	uri := util.GetPath(url)
	span, spanCtx := opentracing.StartSpanFromContext(ctx, uri)
	defer span.Finish()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return spanCtx, nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, method)
	err = span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	if err != nil {
		tracinglogger.ContextLog(spanCtx).Errorf("inject to HTTPHeaders err %v", err)
	}
	resp, err = c.Do(req)
	if err != nil {
		ext.Error.Set(span, true)
		return spanCtx, resp, err
	} else {
		ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))
	}
	return spanCtx, resp, nil
}
