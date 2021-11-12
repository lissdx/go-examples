package main

import (
	"context"
	"time"
)


func Stream(ctx context.Context, out chan<- Value) error {
	// Create a derived Context with a 10s timeout; dctx
	// will be cancelled upon timeout, but ctx will not.
	// cancel is a function that will explicitly cancel dctx.
	dctx, cancel := context.WithTimeout(ctx, time.Second * 10)

	// Release resources if SlowOperation completes before timeout
	defer cancel()

	res, err := SlowOperation(dctx)
	if err != nil {                     // True if dctx times out
		return err
	}

	for {
		select {
		case out <- res:                // Read from res; send to out

		case <-ctx.Done():              // Triggered if ctx is cancelled
			return ctx.Err()
		}
	}
}