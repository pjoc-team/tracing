package main

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"github.com/pjoc-team/tracing/tracingmq"
	"runtime"
	"time"
)

func main() {
	tracing.InitOnlyTracingLog("mq_consumer")
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second
	c, err := nsq.NewConsumer("push_test", "mq_consumer", cfg)
	if err != nil {
		fmt.Println(err)
	}
	c.SetLogger(nil, 0)
	c.AddConcurrentHandlers(&NsqConsumer{}, runtime.NumCPU())
	if err := c.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
		panic(err)
	}
	select {}
}

type NsqConsumer struct{}

func (c *NsqConsumer) HandleMessage(msg *nsq.Message) error {
	fmt.Println("start pull")
	data := &tracingmq.TraceMqData{}
	err := json.Unmarshal(msg.Body, data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	log := tracinglogger.ContextLog(tracing.BuildContextByCarrier(data.Carriers, "mq_consumer", data.Topic))
	log.Infof("pull carriers:%v topic:%s data:%s", data.Carriers, data.Topic, data.Data)
	log.Info("===================")

	return nil
}
