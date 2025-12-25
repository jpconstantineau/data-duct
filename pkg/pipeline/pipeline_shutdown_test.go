package pipeline

import (
	"context"
	"testing"
	"time"
)

func TestPipelineShutdownCompletes_NoHang(t *testing.T) {
	t.Parallel()

	src := func(ctx context.Context) (<-chan int, error) {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for i := 0; i < 1000; i++ {
				select {
				case <-ctx.Done():
					return
				case ch <- i:
				}
			}
		}()
		return ch, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	runnable := New("shutdown", src).
		Then(func(ctx context.Context, n int) (int, error) {
			// Small delay encourages concurrency paths.
			time.Sleep(1 * time.Millisecond)
			return n, nil
		}).
		To(func(ctx context.Context, n int) error { return nil })

	_, _ = runnable.Run(ctx)
	// If Run returns, internal waitgroups completed (no hang).
}
