package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/pjoc-team/tracing/logger"
	"github.com/pjoc-team/tracing/tracing"
	"sync"
)

type nsqout struct {
	Topic    string
	Host     string
	producer *nsq.Producer
	lock     *sync.Once
}

func (nw *nsqout) Write(p []byte) (n int, err error) {
	if nw.producer == nil {
		nw.lock.Do(func() {
			nw.producer, err = nsq.NewProducer(nw.Host, nsq.NewConfig())
		})
	}
	param := map[string]interface{}{}
	//可以追加一些自己的json key
	param["key1"] = "value1"
	json.Unmarshal(p, &param)
	body, _ := json.Marshal(param)
	//上报nsq，方便kibana查看
	//上报log时，注意key的格式，数据类型，是否会跟之前的业务有冲突
	//注意用异步API，防阻塞
	err = nw.producer.PublishAsync(nw.Topic, body, nil)
	//输出到本地 取决于自己的业务需求自行实现
	fmt.Println(string(body))
	return len(p), err
}

func main() {
	tracing.InitOnlyTracingLog("output-nsq")
	logger.SetFormatter(logger.FormatJson)
	logger.SetOutput(&nsqout{
		Topic: "test",
		Host:  "127.0.0.1:4150",
		lock:  &sync.Once{},
	})
	for i := 0; i < 10; i++ {
		tracing.HandleFunc(func(ctx context.Context) {
			log := logger.ContextLog(ctx)
			log.Infof("output-nsq %d", i)
		})
	}
}
