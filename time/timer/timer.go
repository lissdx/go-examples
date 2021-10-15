package main

import (
	"fmt"
	"github.com/lissdx/yapgo/pkg/pipeline"
	"time"
)

const TimerTime = 1
const FuncTime = 5


func longTimeJob() <-chan interface{} {
	outStream := make(chan interface{})
	go func() {
		time.Sleep(time.Second * FuncTime)
		outStream <- FuncTime
	}()
	return outStream
}


func main() {
	timer := time.NewTimer(time.Second * TimerTime)

	done := make(chan interface{})
	defer close(done)

	select {
	case <-timer.C:
		fmt.Printf("New Timer timeout happend")
		done <- struct {}{}
	case v := <-pipeline.OrDone(done, longTimeJob()):
		fmt.Printf("longTimeJob happend: %v", v)
	}

	timer.Stop()
	select {
	case <-timer.C:
	default:
	}

}
