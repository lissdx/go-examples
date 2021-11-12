package cloud_native_patterns

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	isServiceReachable := func() bool {
		d := consecutiveFailures - int(failureThreshold)

		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				//m.RUnlock()
				//return "", errors.New("service unreachable")
				return false
			}
		}

		return true
	}

	return func(ctx context.Context) (string, error) {
		m.RLock()                       // Establish a "read lock"

		//d := consecutiveFailures - int(failureThreshold)
		//
		//if d >= 0 {
		//	shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
		//	if !time.Now().After(shouldRetryAt) {
		//		m.RUnlock()
		//		return "", errors.New("service unreachable")
		//	}
		//}
		if !isServiceReachable(){
			defer m.RUnlock()
			return "", errors.New("service unreachable")
		}

		m.RUnlock()                     // Release read lock

		response, err := circuit(ctx)   // Issue request proper

		m.Lock()                        // Lock around shared resources
		defer m.Unlock()

		lastAttempt = time.Now()        // Record time of attempt

		if err != nil {                 // Circuit returned an error,
			consecutiveFailures++       // so we count the failure
			return response, err        // and return
		}

		consecutiveFailures = 0         // Reset failures counter

		return response, nil
	}
}

func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
	var threshold time.Time
	var result string
	var err error
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		m.Lock()

		defer func() {
			threshold = time.Now().Add(d)
			m.Unlock()
		}()

		if time.Now().Before(threshold) {
			return result, err
		}

		result, err = circuit(ctx)

		return result, err
	}
}