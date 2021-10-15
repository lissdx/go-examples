package main

import (
	"context"
	"fmt"
	"github.com/lissdx/yapgo/pkg/pipeline"
	_ "github.com/lissdx/yapgo/pkg/pipeline"
	"math/rand"
	"strconv"
	"sync"
	_ "sync"
	"time"
)

func worker(ctx context.Context, workerName string, wg *sync.WaitGroup, done <-chan interface{}, out chan<- interface{}) {
	//func worker(ctx context.Context, workerName string, out chan <- interface{})  {
	waitTime := time.Duration(rand.Intn(100)+10) * time.Millisecond
	fmt.Printf("worker %s sleep %d\n", workerName, waitTime)
	timer := time.NewTimer(waitTime)
	defer wg.Done()
	defer func() {
		timer.Stop()
		select {
		case <-timer.C:
			fmt.Printf("drain timer of worker: %s\n", workerName)
		default:
		}
	}()

	select {
	case <-ctx.Done():
		return
	case <-done:

	case <-timer.C:
		select {
		case <-done:
			return
		case out <- fmt.Sprintf("worker %s finished\n", workerName):
		}
	}
}

// Example has a huge issues with writing to closed channel
// Don't pass channel to "write to" to function
func main() {
	var wg sync.WaitGroup
	ctx, cancelFn := context.WithCancel(context.Background())
	outStream := make(chan interface{})
	done := make(chan interface{})
	defer close(done)

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go worker(ctx, strconv.Itoa(i), &wg, done, outStream)
		//go worker(ctx, strconv.Itoa(i), outStream)
	}

	res := <-outStream
	cancelFn()
	wg.Wait()
	close(outStream)
	for v := range pipeline.OrDone(done, outStream) {
		fmt.Printf("DRAIN %s", v)
	}

	fmt.Printf("*** got first msg: %s", res)
	time.Sleep(time.Second * 1)
}
