package pipeline

import (
	"context"
	"testing"
	"time"
)

func TestPipelinePanicBecomesError(t *testing.T) {
	t.Parallel()

	src := func(ctx context.Context) (<-chan int, error) {
		ch := make(chan int, 1)
		ch <- 1
		close(ch)
		return ch, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := New("panic", src).
		Then(func(ctx context.Context, n int) (int, error) {
			panic("nope")
		}).
		To(func(ctx context.Context, n int) error { return nil }).
		Run(ctx)

	if err == nil {
		t.Fatalf("expected error")
	}
}
