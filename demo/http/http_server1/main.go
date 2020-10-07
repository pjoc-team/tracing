package main

import (
	"github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/tracinghttp"
	"io/ioutil"
	"net/http"
)

func main() {
	err := tracing.InitOnlyTracingLog("http_server1")
	if err != nil {
		tracinglogger.Log().Fatal(err)
	}
	tracinglogger.SetLevel(tracinglogger.DebugLevel)
	tracinglogger.SetReportCallerLevel(tracinglogger.InfoLevel)
	tracinghttp.HandleFunc("/sayHello", func(w http.ResponseWriter, r *http.Request) {
		log := tracinglogger.ContextLog(r.Context())
		//log.Infof("uber-trace-id:%s", r.Header.Get("uber-trace-id"))
		_, res, _ := tracinghttp.Get(r.Context(), http.DefaultClient, "http://localhost:8083/sayHello")
		result, _ := ioutil.ReadAll(res.Body)
		log.Debug("res:%s", result)
		res.Body.Close()
	})
	tracinghttp.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		log := tracinglogger.ContextLog(r.Context())
		//log.Infof("uber-trace-id:%s", r.Header.Get("uber-trace-id"))
		result, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		log.Infof("%s", result)
		w.Write([]byte("ok"))
	})
	tracinglogger.Log().Fatal(http.ListenAndServe(":8082", nil))
}
