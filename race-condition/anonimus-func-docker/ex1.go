package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 17; i <=21; i++{
		go func(i int) {
			go func(i int) {
				apiVersion := fmt.Sprintf("v1.%d", i)
				fmt.Printf("apiVersion: %s\n", apiVersion)
			}(i)
		}(i)
	}

	time.Sleep(time.Second * 2)
}

