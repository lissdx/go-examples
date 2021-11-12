package main

import (
	"context"
	"fmt"
	"github.com/lissdx/go-examples/pkg/cloud_native_patterns"
	"net/http"
	"time"
)


const contextKeyName = "urlRequest"

func circuitFun(ctx context.Context) (string, error) {
		//context.WithValue()
		urlRequest := ctx.Value(contextKeyName).(string)
		client := &http.Client{}

		// Create new request
		req, err := http.NewRequest("GET", urlRequest, nil)

		if err != nil {
			return "", err
		}
		req.Header.Set("User-Agent",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/605.1.15 (KHTML, like Gecko")

		resp, err := client.Do(req)
		if err != nil{
			return "", err
		}

		return resp.Status, err
}


func main() {
	//ctx := context.WithValue(context.Background(), contextKeyName, "http://rollingstone.com/ads.txt")
	ctx := context.WithValue(context.Background(), contextKeyName, "http://test.unknown/ads.txt")

	worker := cloud_native_patterns.Breaker(circuitFun, 3)

	for i := 0; i < 20; i++{
		s, err := worker(ctx)
		if err != nil {
			if err.Error() == "service unreachable"{
				time.Sleep(time.Millisecond * 500)
			}
			fmt.Printf("%v\n", err)
			continue
		}
		fmt.Printf("string: %v\n", s)
	}

}
