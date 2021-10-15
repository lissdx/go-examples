package main

import (
	"errors"
	"fmt"
	"github.com/lissdx/yapgo/pkg/pipeline"
	"net/http"
)

type ResponseError struct {
	response *http.Response
	err error
}

type NaiveRequest struct{}

func (NaiveRequest) DoRequest() pipeline.ProcessFn {
	return func(inObj interface{}) interface{} {
		urlRequest, ok := inObj.(string)
		if !ok {
			return ResponseError{nil, errors.New(fmt.Sprintf("cant make string assertion for %v", inObj))}
		}
		client := &http.Client{}

		// Create new request
		req, err := http.NewRequest("GET", urlRequest, nil)

		if err != nil {
			return ResponseError{nil, err}
		}
		req.Header.Set("User-Agent",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/605.1.15 (KHTML, like Gecko")

		resp, err := client.Do(req)
		return ResponseError{resp, err}
	}
}


const TimerWaitMillisecond = 2000

func main() {
	naiveRequest := NaiveRequest{}
	done := make(chan interface{})
	urls := []interface{}{
		"http://rollingstone.com/ads.txt",
		"http://nypost.com/ads.txt",
		"http://ign.com/ads.txt",
		"http://newsweek.com/ads.txt",
		"http://tmi.maariv.co.il/ads.txt",
		"http://justjared.com/ads.txt",
		"http://wordplays.com/ads.txt",
		"http://dailywire.com/ads.txt",
		"http://iflscience.com/ads.txt",
		"http://pcmag.com/ads.txt",
		"http://usmagazine.com/ads.txt",
		"http://tasteofhome.com/ads.txt",
		"http://huffpost.com/ads.txt",
		"http://stylecaster.com/ads.txt",
		"http://billboard.com/ads.txt",
		"http://pleated-jeans.com/ads.txt",
		"http://closerweekly.com/ads.txt",
		"http://elnuevodia.com/ads.txt",
		"http://boston.com/ads.txt",
		"http://complex.com/ads.txt",
	}
	//randFn := func() interface{} { return rand.Intn(50000000) }
	//primerFinder := NaivePrimerFinder{}
	pLine := pipeline.New()
	pLine.AddStageWithFanOut(naiveRequest.DoRequest(), uint64(len(urls)))
	outStream := pLine.Run(done, pipeline.Generator(done, urls...))

	for s := range pipeline.OrDone(done, outStream) {
		reqErr := s.(ResponseError)
		fmt.Printf("Status: %s\n",  reqErr.response.Status)
	}
}
