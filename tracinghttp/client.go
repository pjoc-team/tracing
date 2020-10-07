package tracinghttp

import (
	"context"
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
	client := NewClient(c)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return ctx, nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	resp, err = client.Do(req)
	return ctx, resp, err
}
