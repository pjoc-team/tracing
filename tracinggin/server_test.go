package tracinggin

import (
	"github.com/gin-gonic/gin"
	tracinglogger "github.com/pjoc-team/tracing/logger"
	"net/http"
	"testing"
)

func ExampleTracingServerInterceptor() {
	router := gin.Default()
	router.Use(TracingServerInterceptor())
	router.GET("/test", func(c *gin.Context) {
		//do something
		c.String(http.StatusOK, "gin ok")
	})
	err := router.Run(":8080")
	if err != nil {
		tracinglogger.Log().Fatal(err)
	}
}

func ExampleTracingLogInterceptor() {
	router := gin.New()
	router.Use(
		TracingServerInterceptor(),
		TracingLogInterceptor(
			PanicStack(3, 6),                     /*可选**/
			SkipPath("/test/*"),                  /*不设置则全部不跳过**/
			IncludePath("/info/*"),               /*不设置则包含全部**/
			LogPrefix("gin-status", "gin-panic"), /*可选**/
		))
	router.GET("/test", func(c *gin.Context) {
		//do something
		c.String(http.StatusOK, "gin ok")
	})
	router.GET("/info/show", func(c *gin.Context) {
		//do something
		c.String(http.StatusOK, "gin ok")
	})
	err := router.Run(":8080")
	if err != nil {
		tracinglogger.Log().Fatal(err)
	}
}

func TestSkipPath(t *testing.T) {
	//需要skip
	if path(logOptions{
		skipPaths: []string{"/*"},
	}, "/") {
		panic("not skipPaths")
	}
	//需要skip
	if path(logOptions{
		skipPaths: []string{"/*/info/*"},
	}, "/test/info/get") {
		panic("not skipPaths")
	}
	//不需要skip
	if !path(logOptions{
		skipPaths: []string{"/*/info/*/save"},
	}, "/test/info/get") {
		panic("skipPaths")
	}
}

func TestIncludePaths(t *testing.T) {
	//需要include
	if !path(logOptions{
		includePaths: []string{"/*"},
	}, "/") {
		panic("not includePaths")
	}
	//需要include
	if !path(logOptions{
		includePaths: []string{"/*/info/*"},
	}, "/test/info/get") {
		panic("not includePaths")
	}

	//需要include
	if !path(logOptions{
		includePaths: []string{"/test/info/get/"},
	}, "/test/info/get") {
		panic("not includePaths")
	}

	//不需要include
	if path(logOptions{
		includePaths: []string{"/test/info/get/"},
	}, "/test/info/get/all") {
		panic("includePaths")
	}

	//不需要include
	if path(logOptions{
		includePaths: []string{"/*/info/*/save", "/test/info"},
	}, "/test/info/get") {
		panic("includePaths")
	}
}

func TestPaths(t *testing.T) {
	//不需要include
	if path(logOptions{
		includePaths: []string{"/*/info/*"},
		skipPaths:    []string{"/*/test/*"},
	}, "/") {
		panic("includePaths")
	}

	//需要include
	if !path(logOptions{
		includePaths: []string{"/*/info/*"},
		skipPaths:    []string{"/*/test/*"},
	}, "/user/info/winter") {
		panic("not includePaths")
	}

	//不需要include
	if path(logOptions{
		includePaths: []string{"/*/info/*"},
		skipPaths:    []string{"/*/test/*"},
	}, "/user/test/winter") {
		panic("not includePaths")
	}
}
