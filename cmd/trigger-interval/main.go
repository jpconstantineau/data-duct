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
		interval = flag.Duration("interval", 1*time.Second, "trigger interval")
		duration = flag.Duration("duration", 10*time.Second, "total runtime; 0 means run until interrupted")
	)
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	coalescer := examplesupport.NewCoalescer()

	process := examplesupport.ProcessStep{
		Process: func(ctx context.Context, req examplesupport.IngestionRequest) (examplesupport.IngestionResult, error) {
			start := time.Now()
			// Simulate work that can exceed the trigger interval, so coalescing is observable.
			time.Sleep(150 * time.Millisecond)
			end := time.Now()
			return examplesupport.IngestionResult{StartedAt: start, FinishedAt: end, Records: 1, Outcome: "succeeded"}, nil
		},
	}

	runnerDone := make(chan error, 1)
	go func() {
		runnerDone <- coalescer.Run(ctx, func(ctx context.Context) error {
			firedAt := time.Now()
			ev := examplesupport.TriggerEvent{Kind: "interval", Occurred: firedAt, SourceRef: "simulated://interval"}

			src := func(ctx context.Context) (<-chan examplesupport.TriggerEvent, error) {
				ch := make(chan examplesupport.TriggerEvent, 1)
				ch <- ev
				close(ch)
				return ch, nil
			}

			sink := func(ctx context.Context, res examplesupport.IngestionResult) error {
				examplesupport.PrintOutcome(os.Stdout, "trigger-interval", res.Outcome, fmt.Sprintf("records=%d", res.Records))
				return nil
			}

			_, err := pipeline.New("trigger-interval", src).
				Then(process.Handle).
				To(sink).
				Run(ctx)
			return err
		})
	}()

	loop := examplesupport.NewIntervalLoop()
	_ = loop.Run(ctx, *interval, *duration, func(time.Time) {
		coalescer.Request()
	})

	cancel()
	_ = <-runnerDone
}
