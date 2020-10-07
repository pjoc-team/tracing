package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/tracinggin"
	"net/http"
)

func main() {
	err := tracing.InitOnlyTracingLog("http_servergin")
	if err != nil {
		tracinglogger.Log().Fatal(err)
	}
	tracinglogger.SetLevel(tracinglogger.DebugLevel)
	tracinglogger.MinReportCallerLevel(tracinglogger.InfoLevel)
	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(
		tracinggin.TracingServerInterceptor(),
		tracinggin.TracingLogInterceptor(tracinggin.SkipPath("/sayHello/*")))
	router.GET("/sayHello", func(c *gin.Context) {
		log := tracinglogger.ContextLog(c.Request.Context())
		log.Infoln("hello:8084")
		log.Debug(">>>>>>sayHello>>>>>>")
		c.String(http.StatusOK, "sayHello")
	})
	router.POST("/info", func(c *gin.Context) {
		log := tracinglogger.ContextLog(c.Request.Context())
		log.Infoln("hello:8084", c.PostForm("info"))
		log.Debug(">>>>>>info>>>>>>")
		c.String(http.StatusOK, "info")
	})
	router.GET("/pc", func(c *gin.Context) {
		panic("Test PanicStack")
	})
	err = router.Run(":8084")
	if err != nil {
		tracinglogger.Log().Fatal(err)
	}
}
