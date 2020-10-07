package main

import (
	"github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/tracinghttp"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	err := tracing.InitOnlyTracingLog("http_server2")
	if err != nil {
		logger.Log().Fatal(err)
	}
	tracinghttp.HandleFunc("/sayHello", func(w http.ResponseWriter, r *http.Request) {
		log := logger.ContextLog(r.Context())
		log.Infoln("hello:8083")

		header := make(map[string]string)
		header[tracing.HttpHeaderKeyXRequestID] = time.Now().Format("20060102150405")
		_, res, _ := tracinghttp.GetDo(r.Context(), http.DefaultClient, "http://localhost:8084/sayHello", header)
		result, _ := ioutil.ReadAll(res.Body)
		log.Infof("res:%s", result)
		res.Body.Close()
		w.Write([]byte("i am ok"))
	})
	logger.Log().Fatal(http.ListenAndServe(":8083", nil))
}
