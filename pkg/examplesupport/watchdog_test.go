package examplesupport

import (
	"context"
	"testing"
	"time"
)

type fakeWDicker struct{ ch chan time.Time }

func (t *fakeWDicker) Chan() <-chan time.Time { return t.ch }
func (t *fakeWDicker) Stop()                  {}

type fakeWDTimer struct{ ch chan time.Time }

func (t *fakeWDTimer) Chan() <-chan time.Time { return t.ch }
func (t *fakeWDTimer) Stop() bool             { return true }

func TestWatchdogLoop_Run_EmitsAlertOnDeadline(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tickCh := make(chan time.Time, 10)
	deadlineCh := make(chan time.Time, 1)

	loop := WatchdogLoop{
		NewTicker: func(d time.Duration) WatchdogTicker { return &fakeWDicker{ch: tickCh} },
		NewTimer:  func(d time.Duration) WatchdogTimer { return &fakeWDTimer{ch: deadlineCh} },
		Now:       func() time.Time { return time.Unix(100, 0).UTC() },
	}

	alerts := make(chan AlertSignal, 1)
	done := make(chan error, 1)
	go func() {
		done <- loop.Run(ctx, 5*time.Second, 1*time.Second, func(sig AlertSignal) {
			alerts <- sig
		})
	}()

	// Simulate a couple of poll ticks, then deadline fires.
	tickCh <- time.Unix(101, 0)
	tickCh <- time.Unix(102, 0)
	deadlineCh <- time.Unix(105, 0)

	select {
	case sig := <-alerts:
		if sig.Reason == "" {
			t.Fatalf("expected reason")
		}
		if sig.At.IsZero() {
			t.Fatalf("expected At")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for alert")
	}

	err := <-done
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
