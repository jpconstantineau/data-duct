package examplesupport

import "context"

// Coalescer coalesces multiple trigger requests into at most one pending run.
//
// Request is non-blocking; if a run is in progress and one pending run is already queued,
// additional requests are dropped.
//
// Usage:
//   c := NewCoalescer()
//   go c.Run(ctx, runOnce)
//   c.Request()
//
// This is intentionally small and stdlib-only for runnable examples.
type Coalescer struct {
	reqs chan struct{}
}

func NewCoalescer() *Coalescer {
	return &Coalescer{reqs: make(chan struct{}, 1)}
}

func (c *Coalescer) Request() {
	if c == nil {
		return
	}
	select {
	case c.reqs <- struct{}{}:
	default:
		// already has one pending
	}
}

func (c *Coalescer) Run(ctx context.Context, runOnce func(context.Context) error) error {
	if c == nil || runOnce == nil {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.reqs:
			if err := runOnce(ctx); err != nil {
				return err
			}
		}
	}
}
