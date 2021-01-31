package main

import (
	"bytes"
	"context"
	"github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/tracinghttp"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	err := tracing.InitOnlyTracingLog("http_client")
	if err != nil {
		logger.Log().Fatal(err)
	}
	logger.SetLevel(logger.DebugLevel)
	logger.SetReportCallerLevel(logger.InfoLevel)
	tracing.HandleFunc(func(ctx context.Context) {
		log := logger.ContextLog(ctx)
		log.Debug("start client")
		header := make(map[string]string)
		header[string(tracing.HttpHeaderKeyXRequestID)] = time.Now().Format("20060102150405")
		sctx, _, _ := tracinghttp.GetDo(ctx, http.DefaultClient, "http://localhost:8082/sayHello", header)
		_, res, err := tracinghttp.Post(ctx, http.DefaultClient, "http://localhost:8082/info", "application/json", bytes.NewReader([]byte("{\"info\":\"hello\"}")))
		if err != nil {
			log.Errorf("%v", err)
		} else {
			result, _ := ioutil.ReadAll(res.Body)
			log.Infof("%s", result)
		}
		header["Content-Type"] = "application/json"
		sctx, res, err = tracinghttp.PostDo(ctx, http.DefaultClient, "http://localhost:8082/info", header, bytes.NewReader([]byte("{\"info\":\"hello trace3\"}")))
		sctxLog := logger.ContextLog(sctx)
		if err != nil {
			sctxLog.Errorf("%v", err)
		} else {
			result, _ := ioutil.ReadAll(res.Body)
			sctxLog.Infof("%s", result)
		}
		tracing.Forward(sctx, "testForward", func(ctx context.Context) {
			log = logger.ContextLog(ctx)
			log.Println("test testForward")
		})

		tracinghttp.Post(ctx, http.DefaultClient, "http://localhost:8084/info", "application/x-www-form-urlencoded", bytes.NewReader([]byte("info=jaeger")))

		tracinghttp.GetDo(ctx, http.DefaultClient, "http://localhost:8089/index", header)
		tracinghttp.GetDo(ctx, http.DefaultClient, "http://localhost:8084/pc", header)
		log.Info("finish client")
	})

	//调试上报jaeger ui ，可以打开注释，因为span是异步上报，进程结束可能会导致未执行udp 上报jaeger agent
	//select {}
}
