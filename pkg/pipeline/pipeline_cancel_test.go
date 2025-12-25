package pipeline

import (
	"context"
	"testing"
	"time"
)

func TestPipelineCancellationReturnsPromptly(t *testing.T) {
	t.Parallel()

	src := func(ctx context.Context) (<-chan int, error) {
		ch := make(chan int)
		go func() {
			defer close(ch)
			// Emit until cancelled.
			for i := 0; ; i++ {
				select {
				case <-ctx.Done():
					return
				case ch <- i:
				}
			}
		}()
		return ch, nil
	}

	sink := func(ctx context.Context, s string) error {
		// Simulate work.
		time.Sleep(2 * time.Millisecond)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	runnable := New("cancel", src).
		Then(func(ctx context.Context, n int) (string, error) { return itoa(n), nil }).
		To(sink)

	done := make(chan struct{})
	var res Result
	var err error
	go func() {
		defer close(done)
		res, err = runnable.Run(ctx)
	}()

	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("expected Run to return promptly after cancel")
	}

	if err == nil {
		t.Fatalf("expected cancellation error")
	}
	if res.State() != StateCancelled {
		t.Fatalf("expected cancelled, got %s", res.State())
	}
}
