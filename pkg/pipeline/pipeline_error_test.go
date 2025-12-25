package pipeline

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestPipelineProcessorErrorPropagates(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("boom")
	src := func(ctx context.Context) (<-chan int, error) {
		ch := make(chan int, 3)
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)
		return ch, nil
	}

	sink := func(ctx context.Context, s string) error { return nil }

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := New("err", src).
		Then(func(ctx context.Context, n int) (string, error) {
			if n == 2 {
				return "", sentinel
			}
			return itoa(n), nil
		}).
		To(sink).
		Run(ctx)

	if !errors.Is(err, sentinel) {
		t.Fatalf("expected error %v, got %v", sentinel, err)
	}
	if res.State() != StateFailed {
		t.Fatalf("expected failed, got %s", res.State())
	}
}
