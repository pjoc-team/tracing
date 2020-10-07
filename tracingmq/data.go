package tracingmq

import "github.com/opentracing/opentracing-go"

type TraceMqData struct {
	Carriers opentracing.TextMapCarrier
	Topic    string
	Data     []byte
}
