package examplesupport

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestCoalescer_CoalescesToOnePendingRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewCoalescer()

	blockFirstRun := make(chan struct{})
	releaseFirstRun := make(chan struct{})
	var runs atomic.Int32

	done := make(chan error, 1)
	go func() {
		err := c.Run(ctx, func(ctx context.Context) error {
			n := runs.Add(1)
			if n == 1 {
				close(blockFirstRun)
				select {
				case <-releaseFirstRun:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
		done <- err
	}()

	c.Request()

	select {
	case <-blockFirstRun:
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for first run to start")
	}

	// While first run is blocked, spam requests; only one should be pending.
	for i := 0; i < 100; i++ {
		c.Request()
	}

	close(releaseFirstRun)

	// Wait until the second run happens.
	deadline := time.Now().Add(2 * time.Second)
	for runs.Load() < 2 {
		if time.Now().After(deadline) {
			t.Fatalf("expected 2 runs (1 + 1 pending), got %d", runs.Load())
		}
		time.Sleep(5 * time.Millisecond)
	}

	cancel()
	_ = <-done

	if got := runs.Load(); got != 2 {
		t.Fatalf("expected exactly 2 runs, got %d", got)
	}
}
