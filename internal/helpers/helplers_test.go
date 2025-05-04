package helpers

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var errRand = errors.New("temporary failure")

func TestRetry(t *testing.T) {
	dsu := func() error {
		if rand.Intn(11) < 9 {
			fmt.Println("Operation failed, retrying...")
			return errRand
		}
		fmt.Println("Operation succeeded!")
		return nil
	}
	err := WithRetry(context.Background(), 10, 5*time.Second, dsu)
	fmt.Printf("%v\n", err)
	if err != nil && !errors.Is(err, errRand) {
		t.Errorf(`WithRetry() = %v, want match for %v, ErrorIs = %v`, err, errRand, errors.Is(err, errRand))
	}
}
