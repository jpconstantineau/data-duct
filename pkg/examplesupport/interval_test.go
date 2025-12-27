package examplesupport

import (
	"context"
	"testing"
	"time"
)

type fakeTicker struct{ ch chan time.Time }

func (t *fakeTicker) Chan() <-chan time.Time { return t.ch }
func (t *fakeTicker) Stop()                  {}

type fakeTimer struct{ ch chan time.Time }

func (t *fakeTimer) Chan() <-chan time.Time { return t.ch }
func (t *fakeTimer) Stop() bool             { return true }

func TestIntervalLoop_Run_TicksUntilDuration(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tickCh := make(chan time.Time, 10)
	timerCh := make(chan time.Time)

	loop := IntervalLoop{
		NewTicker: func(d time.Duration) Ticker { return &fakeTicker{ch: tickCh} },
		NewTimer:  func(d time.Duration) Timer { return &fakeTimer{ch: timerCh} },
	}

	var got int
	twoTicks := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- loop.Run(ctx, 10*time.Millisecond, 100*time.Millisecond, func(time.Time) {
			got++
			if got == 2 {
				close(twoTicks)
			}
		})
	}()

	// Two ticks, then duration expires.
	tickCh <- time.Unix(1, 0)
	tickCh <- time.Unix(2, 0)

	select {
	case <-twoTicks:
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for two ticks")
	}

	go func() { timerCh <- time.Unix(3, 0) }()

	err := <-done
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 2 {
		t.Fatalf("expected 2 ticks, got %d", got)
	}
}
