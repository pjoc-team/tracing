package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/tracingmq"
	"strconv"
)

func main() {
	tracing.InitOnlyTracingLog("mq_producer")
	producer, err := nsq.NewProducer("127.0.0.1:4150", nsq.NewConfig())
	if err != nil {
		fmt.Println(err)
	} else {
		for i := 0; i < 10; i++ {
			tracing.HandleFunc(func(ctx context.Context) {
				log := tracinglogger.ContextLog(ctx)
				log.Infof("start push_test %d", i)
				tarceMqData := tracingmq.TracingMqProducer(ctx, "push_test", []byte(strconv.Itoa(i)))
				byte, _ := json.Marshal(tarceMqData)
				producer.Publish("push_test", byte)
				if err != nil {
					log.Error(err)
				}
				log.Infof("finish push_test %d", i)
			})
		}
	}
}
