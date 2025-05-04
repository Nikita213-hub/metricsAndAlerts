package helpers

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type RetryError struct {
	retryCount int
	err        error
}

func NewRetryError(rc int, err error) error {
	return &RetryError{
		retryCount: rc,
		err:        err,
	}
}

func (re *RetryError) Error() string {
	return fmt.Sprintf("Max count of retries(%d): %v", re.retryCount, re.err)
}

func (re *RetryError) Unwrap() error {
	return re.err
}

func WithRetry(ctx context.Context, maxTries int, maxDelay time.Duration, callable func() error) error {
	rc := 0
	var err error
	baseTime := 1 * time.Second
	for retries := 0; retries < maxTries; retries++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		rc++
		err = callable()
		if err == nil {
			return nil
		}
		backoffTime := baseTime * time.Duration(math.Pow(2, float64(retries)))
		jitter := time.Duration(rand.Int63n(int64(backoffTime)))
		if jitter > maxDelay {
			jitter = maxDelay
		}
		time.Sleep(jitter)
	}
	return NewRetryError(rc, err)
}
