package tracinghttp

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := &http.Client{Timeout: 500 * time.Millisecond}
	// 使用tracinghttp，上报耗时
	client = NewClient(client)
	// 这里的context建议使用tracing context，这样可以把上下文串起来
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://baidu.com", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	get, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(get)

}
