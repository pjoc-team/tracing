package tracing

import (
	"fmt"
	"os"
	"testing"
)

func ExampleInitOnlyTracingLog() {
	//必须指定serviceName
	//只生成tracing信息的日志，不会上报jaeger服务器
	err := InitOnlyTracingLog("demo")
	//初始化错误要记得处理，方便排查问题
	if err != nil {
		panic(err)
	}
}

func ExampleInitFromConfig() {
	err := InitFromConfig(&Config{
		ServiceName: "demo",
		JaegerAgent: "127.0.0.1:6831",
		//rate<=0永不采样  rate<1按固定概率采样  rate=1总是采样  rate>1限制每秒最大采样数量
		SamplerRateLimit: 0.5,
		//缓冲区队列大小，满了后续span则丢弃
		QueueSize: 500,
		//缓冲区刷新间隔，单位s
		BufferFlushInterval: 1,
	})
	//初始化错误要记得处理，方便排查问题
	if err != nil {
		panic(err)
	}
}

func ExampleInitFromEnv() {
	//配置以下环境变量，各个参数根据应用实际情况调整
	/*
	   - name: JAEGER_SERVICE_NAME
	     value: "demo"
	   - name: JAEGER_AGENT_HOST
	     value: "127.0.0.1"
	   - name: JAEGER_AGENT_PORT
	     value: "6831"
	   - name: JAEGER_REPORTER_LOG_SPANS
	     value: "true"
	   - name: JAEGER_REPORTER_MAX_QUEUE_SIZE
	     value: "500"
	   - name: JAEGER_REPORTER_FLUSH_INTERVAL
	     value: "1s"
	   - name: JAEGER_SAMPLER_TYPE
	     value: "probabilistic"
	   - name: JAEGER_SAMPLER_PARAM
	     value: "1"
	*/
	err := InitFromEnv()
	//初始化错误要记得处理，方便排查问题
	if err != nil {
		panic(err)
	}
}

func TestInitFromEnv(t *testing.T) {
	os.Setenv("JAEGER_SERVICE_NAME", "JaegerTest")
	os.Setenv("JAEGER_AGENT_HOST", "localhost")
	os.Setenv("JAEGER_AGENT_PORT", "6831")
	os.Setenv("JAEGER_REPORTER_LOG_SPANS", "true")
	os.Setenv("JAEGER_REPORTER_MAX_QUEUE_SIZE", "500")
	os.Setenv("JAEGER_REPORTER_FLUSH_INTERVAL", "1s")
	os.Setenv("JAEGER_SAMPLER_TYPE", "probabilistic")
	os.Setenv("JAEGER_SAMPLER_PARAM", "1")
	err := InitFromEnv()
	if err != nil {
		panic(err)
	}
}

func TestInitOnlyTracingLog(t *testing.T) {
	err := InitOnlyTracingLog("JaegerTest")
	if err != nil {
		panic(err)
	}
}

func TestInitFromConfig(t *testing.T) {
	c := &Config{
		ServiceName:      "demo",
		SamplerRateLimit: 1,
	}
	err := InitFromConfig(c)
	if err != nil {
		panic(err)
	}
}

func TestInitOnceErr(t *testing.T) {
	c := &Config{
		SamplerRateLimit: 1,
		JaegerAgent:      "127.0.0.1:6831",
	}
	err := InitFromConfig(c)
	if err != nil {
		fmt.Println(">>>>>>InitFromConfig 1 err ", err)
	}
	err = InitFromEnv()
	if err != nil {
		fmt.Println(">>>>>>InitFromEnv 2 err ", err)
	}
	err = InitOnlyTracingLog("")
	if err != nil {
		fmt.Println(">>>>>>InitFromEnv 3 err ", err)
	}
}

func TestInitOnceSuccess(t *testing.T) {
	c := &Config{
		ServiceName:      "demo",
		SamplerRateLimit: 1,
		JaegerAgent:      "127.0.0.1:6831",
	}
	err := InitFromConfig(c)
	if err != nil {
		fmt.Println(">>>>>>InitFromConfig 1 err ", err)
	}
	err = InitFromEnv()
	if err != nil {
		fmt.Println(">>>>>>InitFromEnv 2 err ", err)
	}
	err = InitOnlyTracingLog("")
	if err != nil {
		fmt.Println(">>>>>>InitFromEnv 3 err ", err)
	}
}
