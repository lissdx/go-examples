package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	waitChan := make(chan bool)
	syncChan := make(chan int)

	go func() {
		var i int = 1

		for true {
			val := <- syncChan
			if val + 1 == i {
				fmt.Printf("%d\n", i)
				syncChan <- i
				i += 3
				time.Sleep(time.Millisecond * 1)
				continue
			} else {
				syncChan <- val
				runtime.Gosched()
			}
		}
	}()

	go func() {
		var i int = 2

		for true {
			val := <- syncChan
			if val + 1 == i {
				fmt.Printf("%d\n", i)
				syncChan <- i
				i += 3
				time.Sleep(time.Millisecond * 1)
				continue
			} else {
				syncChan <- val
				runtime.Gosched()
			}
		}
	}()

	go func() {
		var i int = 3

		for true {
			val := <- syncChan
			if val + 1 == i {
				fmt.Printf("%d\n", i)
				syncChan <- i
				i += 3
				time.Sleep(time.Millisecond * 1)
				continue
			} else {
				syncChan <- val
				runtime.Gosched()
			}
		}
	}()

	syncChan <- 0
	<-waitChan
}
