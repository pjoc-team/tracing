package tracinghttp

import (
	tracinglogger "github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"net/http"
)

func ExampleHandleFunc() {
	HandleFunc("/sayHello", func(w http.ResponseWriter, r *http.Request) {
		//do something
	})
	http.ListenAndServe(":8084", nil)
}

func ExampleTracingServerInterceptor() {
	tracing.InitOnlyTracingLog("test")
	mux := http.NewServeMux()
	mux.Handle("/index", &indexHandler{})
	s := &http.Server{
		Addr:    ":8089",
		Handler: TracingServerInterceptor(mux),
	}
	tracinglogger.Log().Fatal(s.ListenAndServe())
}

type indexHandler struct{}

func (ih *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := tracinglogger.ContextLog(r.Context())
	log.Infoln("hello:8089")
	w.Write([]byte("Welcome"))
}
