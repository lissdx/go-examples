package main

import (
	"fmt"
	"github.com/lissdx/yapgo/pkg/pipeline"
	"math/rand"
	"time"
)

type Primer struct {
	IsPrimer bool
	integer  int
}

type NaivePrimerFinder struct{}

func (NaivePrimerFinder) isPrime(integer int) bool {
	isPrime := true
	for divisor := integer - 1; divisor > 1; divisor-- {
		if integer%divisor == 0 {
			isPrime = false
			break
		}
	}
	return isPrime
}

func (n NaivePrimerFinder) NaivePrimer() pipeline.ProcessFn {
	return func(inObj interface{}) interface{} {
		intVal, ok := inObj.(int)
		if !ok {
			return Primer{integer: intVal, IsPrimer: false}
		}
		return Primer{integer: intVal, IsPrimer: n.isPrime(intVal)}
	}
}

func (NaivePrimerFinder) endStubFn() pipeline.ProcessFn {
	return func(inObj interface{}) interface{} {
		fmt.Printf("NaivePrimerFinder is pimer: %v\n", inObj)
		return inObj
	}
}

const TimerWaitMillisecond = 2000

func main() {
	randFn := func() interface{} { return rand.Intn(50000000) }
	primerFinder := NaivePrimerFinder{}
	pLine := pipeline.New()
	timerWait := time.NewTimer(time.Millisecond * TimerWaitMillisecond)
	done := make(chan interface{})

	// Some time we could get:
	// pipeline finished
	// Search took: 1.878421909s
	// timerWait.C 2021-10-09 16:54:08.128891 +0300 IDT m=+2.004311220
	defer func() {
		time.Sleep(time.Second * 1)
		timerWait.Stop()
		select {
		case o := <-timerWait.C:
			fmt.Printf("timerWait.C %v\n", o)
		default:
		}
	}()
	defer close(done)

	pLine.AddStageWithFanOut(primerFinder.NaivePrimer(), 10)
	pLine.AddStage(primerFinder.endStubFn())

	intStream := pipeline.Take(done, pipeline.RepeatFn(done, randFn), 100)

	start := time.Now()
	doneCh := pLine.RunPlug(done, intStream) // will be closed by pLine.RunPlug

	select{
		case <-doneCh:
			fmt.Printf("pipeline finished\n")
	case <- timerWait.C:
		fmt.Printf("timer fired\n")
	}


	fmt.Printf("Search took: %v\n", time.Since(start))
}
