package main

import (
	"github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/tracinghttp"
	"log"
	"net/http"
)

func main() {
	tracing.InitOnlyTracingLog("http_server")
	tracinghttp.HandleFunc("/sayHello", func(w http.ResponseWriter, r *http.Request) {
		log := tracinglogger.ContextLog(r.Context())
		log.Infoln("hello:8084")
		w.Write([]byte("i am ok"))
	})
	log.Fatal(http.ListenAndServe(":8084", nil))
}
