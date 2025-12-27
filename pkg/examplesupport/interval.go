package examplesupport

import (
	"context"
	"time"
)

// Ticker abstracts time.Ticker for deterministic tests.
type Ticker interface {
	Chan() <-chan time.Time
	Stop()
}

// Timer abstracts time.Timer for deterministic tests.
type Timer interface {
	Chan() <-chan time.Time
	Stop() bool
}

type realTicker struct{ t *time.Ticker }

func (t *realTicker) Chan() <-chan time.Time { return t.t.C }
func (t *realTicker) Stop()                  { t.t.Stop() }

type realTimer struct{ t *time.Timer }

func (t *realTimer) Chan() <-chan time.Time { return t.t.C }
func (t *realTimer) Stop() bool             { return t.t.Stop() }

// IntervalLoop runs a callback on each interval tick until duration elapses (if >0)
// or until ctx is cancelled.
type IntervalLoop struct {
	NewTicker func(d time.Duration) Ticker
	NewTimer  func(d time.Duration) Timer
}

func NewIntervalLoop() IntervalLoop {
	return IntervalLoop{
		NewTicker: func(d time.Duration) Ticker { return &realTicker{t: time.NewTicker(d)} },
		NewTimer:  func(d time.Duration) Timer { return &realTimer{t: time.NewTimer(d)} },
	}
}

func (l IntervalLoop) Run(ctx context.Context, interval time.Duration, duration time.Duration, onTick func(time.Time)) error {
	if onTick == nil {
		return nil
	}
	if l.NewTicker == nil {
		l = NewIntervalLoop()
	}

	ticker := l.NewTicker(interval)
	defer ticker.Stop()

	var timer Timer
	if duration > 0 {
		if l.NewTimer == nil {
			l = NewIntervalLoop()
		}
		timer = l.NewTimer(duration)
		defer timer.Stop()
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-ticker.Chan():
			onTick(t)
		case <-func() <-chan time.Time {
			if timer == nil {
				return nil
			}
			return timer.Chan()
		}():
			return nil
		}
	}
}
