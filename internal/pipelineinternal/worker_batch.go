package pipelineinternal

import (
	"context"
	"time"
)

func workerBatch(ctx context.Context, in <-chan feed, out chan<- feed, handler BatchHandler, policy BatchPolicy, logger Logger) {
	defer close(out)

	if policy.Size < 1 {
		policy.Size = 1
	}

	var (
		buf   = make([]feed, 0, policy.Size)
		timer *time.Timer
	)

	resetTimer := func() {
		if policy.MaxWait <= 0 {
			return
		}
		if timer == nil {
			timer = time.NewTimer(policy.MaxWait)
			return
		}
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(policy.MaxWait)
	}

	flush := func() {
		if len(buf) == 0 {
			return
		}

		inputs := make([]any, 0, len(buf))
		for _, f := range buf {
			inputs = append(inputs, f.Data)
		}

		outs, err := handler(ctx, inputs)
		if err == nil {
			for _, o := range outs {
				select {
				case <-ctx.Done():
					// stop emitting
					buf = buf[:0]
					return
				case out <- feed{RootCtx: ctx, PipelineName: buf[0].PipelineName, Data: o}:
				}
			}
		}

		buf = buf[:0]
	}

	for {
		var timerC <-chan time.Time
		if timer != nil {
			timerC = timer.C
		}

		select {
		case <-ctx.Done():
			// Best-effort flush of buffered items on cancel.
			flush()
			return
		case <-timerC:
			flush()
		case f, ok := <-in:
			if !ok {
				flush()
				return
			}
			if len(buf) == 0 {
				resetTimer()
			}
			buf = append(buf, f)
			if len(buf) >= policy.Size {
				flush()
				resetTimer()
			}
		}
	}
}
