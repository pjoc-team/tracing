package main

import (
	"github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/tracinghttp"
	"net/http"
)

type indexHandler struct{}

func (ih *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := logger.ContextLog(r.Context())
	log.Infoln("hello:8089")
	w.Write([]byte("Welcome"))
}

func main() {
	tracing.InitOnlyTracingLog("http_server_mux")
	mux := http.NewServeMux()
	mux.Handle("/index", &indexHandler{})
	s := &http.Server{
		Addr:    ":8089",
		Handler: tracinghttp.TracingServerInterceptor(mux),
	}
	logger.Log().Fatal(s.ListenAndServe())
}
