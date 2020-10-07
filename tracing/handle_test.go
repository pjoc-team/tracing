package tracing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pjoc-team/tracing/util"
	"testing"
)

func ExampleHandleFunc() {
	for i := 0; i < 10; i++ {
		HandleFunc(func(ctx context.Context) {
			//每次都生成trace信息传递下去
			mockPushNsq(ctx)
		})
	}
}

func mockPushNsq(ctx context.Context) {}

func ExampleForward() {
	//假设这是一个服务端的接口
	mockGrpcInterface(func(ctx context.Context) {
		//context.Context 请使用服务端接口传递下来的，不要使用 context.Background()
		Forward(ctx, "getRedis", func(ctx context.Context) {
			//do something
			//需要手动埋点上报的代码块
		})

		Forward(ctx, "saveMysql", func(ctx context.Context) {
			//do something
			//需要手动埋点上报的代码块
		})

	})
}

func mockGrpcInterface(mock func(ctx context.Context)) {
	mock(context.Background())
}

func TestStart(t *testing.T) {
	c := &Config{
		ServiceName:      "JaegerTest",
		SamplerRateLimit: 1,
		JaegerAgent:      "localhost:6381",
	}
	err := InitFromConfig(c)
	if err != nil {
		panic(err)
	}
	span, ctx := Start(context.Background(), "testStart")
	defer span.Finish()
	//do something
	fmt.Println(util.GetTraceInfo(ctx))
}

func ExampleStart() {
	span, ctx := Start(context.Background(), "testStart")
	defer span.Finish()
	//do something
	fmt.Println(util.GetTraceInfo(ctx))
}

func ExampleBuildContextByCarrier() {
	data := opentracing.TextMapCarrier{}
	json.Unmarshal([]byte(mockConsumerMessage()), data)
	BuildContextByCarrier(data, "mq_consumer", "test")
}

func mockConsumerMessage() string {
	return "{}"
}
