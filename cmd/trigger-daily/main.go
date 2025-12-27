package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jpconstantineau/data-duct/pkg/examplesupport"
	"github.com/jpconstantineau/data-duct/pkg/pipeline"
)

func main() {
	var (
		hhmm = flag.String("time", "", "daily time in HH:MM (24h)")
		once = flag.Bool("once", true, "run once then exit")
	)
	flag.Parse()

	now := time.Now()
	var hour, minute int
	if *hhmm == "" {
		// Default to the next minute to make it likely the example runs soon.
		nextMinute := now.Truncate(time.Minute).Add(1 * time.Minute)
		hour, minute = nextMinute.Hour(), nextMinute.Minute()
		*hhmm = examplesupport.FormatHHMM(hour, minute)
	} else {
		h, m, err := examplesupport.ParseHHMM(*hhmm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid -time %q (expected HH:MM)\n", *hhmm)
			os.Exit(2)
		}
		hour, minute = h, m
	}

	target := examplesupport.NextDailyAt(now, hour, minute)
	wait := time.Until(target)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	coalescer := examplesupport.NewCoalescer()

	process := examplesupport.ProcessStep{
		Process: func(ctx context.Context, req examplesupport.IngestionRequest) (examplesupport.IngestionResult, error) {
			start := time.Now()
			time.Sleep(200 * time.Millisecond)
			end := time.Now()
			return examplesupport.IngestionResult{StartedAt: start, FinishedAt: end, Records: 1, Outcome: "succeeded"}, nil
		},
	}

	runnerDone := make(chan error, 1)
	go func() {
		runnerDone <- coalescer.Run(ctx, func(ctx context.Context) error {
			firedAt := time.Now()
			ev := examplesupport.TriggerEvent{Kind: "daily", Occurred: firedAt, SourceRef: "simulated://daily"}

			src := func(ctx context.Context) (<-chan examplesupport.TriggerEvent, error) {
				ch := make(chan examplesupport.TriggerEvent, 1)
				ch <- ev
				close(ch)
				return ch, nil
			}

			sink := func(ctx context.Context, res examplesupport.IngestionResult) error {
				examplesupport.PrintOutcome(os.Stdout, "trigger-daily", res.Outcome, fmt.Sprintf("records=%d", res.Records))
				return nil
			}

			_, err := pipeline.New("trigger-daily", src).
				Then(process.Handle).
				To(sink).
				Run(ctx)
			return err
		})
	}()

	fmt.Fprintf(os.Stdout, "[trigger-daily] configured time=%s next=%s wait=%s\n", *hhmm, examplesupport.FormatTime(target), wait)

	// Wait until the scheduled time.
	timer := time.NewTimer(wait)
	select {
	case <-ctx.Done():
		_ = timer.Stop()
		return
	case <-timer.C:
		coalescer.Request()
	}

	if *once {
		// Give the run a moment to start/finish, then exit.
		time.Sleep(250 * time.Millisecond)
		cancel()
		_ = <-runnerDone
		return
	}

	// If not once, keep triggering every 24h.
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			_ = <-runnerDone
			return
		case <-ticker.C:
			coalescer.Request()
		}
	}
}
