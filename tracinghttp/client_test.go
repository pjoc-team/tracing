package tracinghttp

import (
	"bytes"
	"context"
	"net/http"
)

func ExampleGet() {
	Get(context.Background(), http.DefaultClient, "http://localhost:8084/info")
}

func ExampleGetDo() {
	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	GetDo(context.Background(), http.DefaultClient, "http://localhost:8084/info", header)
}

func ExamplePost() {
	Post(context.Background(), http.DefaultClient, "http://localhost:8084/info", "application/x-www-form-urlencoded", bytes.NewReader([]byte("info=jaeger")))
}

func ExamplePostDo() {
	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	PostDo(context.Background(), http.DefaultClient, "http://localhost:8084/info", header, bytes.NewReader([]byte("info=jaeger")))
}
