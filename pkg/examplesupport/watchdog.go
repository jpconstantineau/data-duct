package examplesupport

import (
	"context"
	"time"
)

// WatchdogTicker abstracts time.Ticker for deterministic tests.
type WatchdogTicker interface {
	Chan() <-chan time.Time
	Stop()
}

// WatchdogTimer abstracts time.Timer for deterministic tests.
type WatchdogTimer interface {
	Chan() <-chan time.Time
	Stop() bool
}

type realWatchdogTicker struct{ t *time.Ticker }

func (t *realWatchdogTicker) Chan() <-chan time.Time { return t.t.C }
func (t *realWatchdogTicker) Stop()                  { t.t.Stop() }

type realWatchdogTimer struct{ t *time.Timer }

func (t *realWatchdogTimer) Chan() <-chan time.Time { return t.t.C }
func (t *realWatchdogTimer) Stop() bool             { return t.t.Stop() }

// WatchdogLoop waits until a deadline elapses and then emits an AlertSignal.
//
// It ticks on poll to make behavior visible (and testable), but it does not
// require any external infrastructure.
type WatchdogLoop struct {
	NewTicker func(d time.Duration) WatchdogTicker
	NewTimer  func(d time.Duration) WatchdogTimer
	Now       func() time.Time
}

func NewWatchdogLoop() WatchdogLoop {
	return WatchdogLoop{
		NewTicker: func(d time.Duration) WatchdogTicker { return &realWatchdogTicker{t: time.NewTicker(d)} },
		NewTimer:  func(d time.Duration) WatchdogTimer { return &realWatchdogTimer{t: time.NewTimer(d)} },
		Now:       time.Now,
	}
}

// Run waits until deadline elapses (unless ctx is canceled), then calls onAlert.
func (l WatchdogLoop) Run(ctx context.Context, deadline time.Duration, poll time.Duration, onAlert func(AlertSignal)) error {
	if onAlert == nil {
		return nil
	}
	if l.NewTicker == nil || l.NewTimer == nil || l.Now == nil {
		l = NewWatchdogLoop()
	}

	if deadline <= 0 {
		onAlert(AlertSignal{At: l.Now(), Reason: "deadline must be >0"})
		return nil
	}
	if poll <= 0 {
		poll = 1 * time.Second
	}

	ticker := l.NewTicker(poll)
	defer ticker.Stop()
	timer := l.NewTimer(deadline)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.Chan():
			// Poll tick: in a real system we'd check last-seen arrivals.
			// Example keeps this hook visible without external dependencies.
		case <-timer.Chan():
			onAlert(AlertSignal{At: l.Now(), Reason: "deadline missed"})
			return nil
		}
	}
}
