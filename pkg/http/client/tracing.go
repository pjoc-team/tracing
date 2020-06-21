package client

import (
	"github.com/jinzhu/copier"
	"net/http"
	"net/http/httptrace"
)

type TracingHttp struct {
	httptrace.ClientTrace
}

func (t *TracingHttp) RoundTrip(request *http.Request) (*http.Response, error) {
}

func NewHttpClient(client *http.Client) *http.Client {
	//req, _ := http.NewRequest("GET", "http://example.com", nil)
	//trace := &httptrace.ClientTrace{
	//	GotConn: func(connInfo httptrace.GotConnInfo) {
	//		fmt.Printf("Got Conn: %+v\n", connInfo)
	//	},
	//	DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
	//		fmt.Printf("DNS Info: %+v\n", dnsInfo)
	//	},
	//}
	//req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	//_, err := http.DefaultTransport.RoundTrip(req)
	//if err != nil {
	//	log.Fatal(err)
	//}
	if client == nil {
		client = http.DefaultClient
	}
	if client.Transport == nil {
		client.Transport = http.DefaultTransport
	}
	t := &TracingHttp{}
	c := &http.Client{
		Transport: t,
	}
	_ = copier.Copy(c, client)

	return c
}
