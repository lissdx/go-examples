package main
import (
	"fmt"
	"time"
)

// https://medium.com/@anandpillai/synchronization-in-go-using-concurrency-primitives-a-case-study-535bb2a71c13
// The problem is — How do you synchronize the
// go-routines so that they print the numbers nicely in order as 1,2,3,4,5 … ?
func main() {
	waitChan := make(chan bool)

	go func()  {
		var i int  = 1

		for true {
			fmt.Printf("%d\n", i)
			i += 3
			time.Sleep(time.Duration(1)*time.Second)
		}
	}()

	go func()  {
		var i int  = 2

		for true {
			fmt.Printf("%d\n", i)
			i += 3
			time.Sleep(time.Duration(1)*time.Second)
		}
	}()

	go func()  {
		var i int  = 3

		for true {
			fmt.Printf("%d\n", i)
			i += 3
			time.Sleep(time.Duration(1)*time.Second)
		}
	}()

	<-waitChan
}