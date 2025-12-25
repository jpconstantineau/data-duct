package pipeline

import (
	"context"
	"testing"
	"time"
)

func TestPipelineBatchFlushOnCancel(t *testing.T) {
	t.Parallel()

	src := func(ctx context.Context) (<-chan int, error) {
		ch := make(chan int)
		go func() {
			defer close(ch)
			// Emit a couple items then wait for cancel.
			for i := 1; i <= 2; i++ {
				select {
				case <-ctx.Done():
					return
				case ch <- i:
				}
			}
			<-ctx.Done()
		}()
		return ch, nil
	}

	flushed := make(chan struct{}, 1)
	handler := func(ctx context.Context, inputs []int) ([]int, error) {
		// If we ever see a partial batch, we consider that a flush.
		if len(inputs) > 0 {
			select {
			case flushed <- struct{}{}:
			default:
			}
		}
		return inputs, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	runnable := New("flush", src).
		ThenBatch(handler, BatchPolicy{Size: 10}).
		To(func(ctx context.Context, n int) error { return nil })

	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = runnable.Run(ctx)
	}()

	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("expected Run to return")
	}

	select {
	case <-flushed:
		// ok
	default:
		t.Fatalf("expected batch to flush buffered items on cancel")
	}
}
