package tracingmq

import (
	"context"
	"encoding/json"
)

func ExampleTracingMqProducer() {
	//比如某个业务需要发送mq消息
	mockMqProducer(func(ctx context.Context) {
		tarceMqData := TracingMqProducer(ctx, "push_test", []byte("test"))
		byte, _ := json.Marshal(tarceMqData)
		mockMqPublish("push_test", byte)
	})
}

func mockMqProducer(mock func(ctx context.Context)) {
	mock(context.Background())
}

func mockMqPublish(topic string, data []byte) {
}
