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
		deadline = flag.Duration("deadline", 5*time.Second, "deadline before alert triggers")
		poll     = flag.Duration("poll", 1*time.Second, "poll interval")
	)
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	loop := examplesupport.NewWatchdogLoop()

	fmt.Fprintf(os.Stdout, "[trigger-watchdog-sla] started deadline=%s poll=%s\n", deadline.String(), poll.String())

	alertCh := make(chan examplesupport.AlertSignal, 1)
	err := loop.Run(ctx, *deadline, *poll, func(sig examplesupport.AlertSignal) {
		fmt.Fprintln(os.Stdout, "[trigger-watchdog-sla] ALERT deadline missed")
		alertCh <- sig
	})
	if err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "[trigger-watchdog-sla] error: %v\n", err)
		os.Exit(1)
	}

	// Emit explicit AlertSignal event into a pipeline.
	sig := <-alertCh
	src := func(ctx context.Context) (<-chan examplesupport.AlertSignal, error) {
		ch := make(chan examplesupport.AlertSignal, 1)
		ch <- sig
		close(ch)
		return ch, nil
	}

	then := func(ctx context.Context, s examplesupport.AlertSignal) (examplesupport.IngestionResult, error) {
		return examplesupport.IngestionResult{StartedAt: s.At, FinishedAt: s.At, Records: 0, Outcome: "ALERT"}, nil
	}

	sink := func(ctx context.Context, res examplesupport.IngestionResult) error {
		examplesupport.PrintOutcome(os.Stdout, "trigger-watchdog-sla", res.Outcome, "deadline missed")
		return nil
	}

	fmt.Fprintln(os.Stdout, "[trigger-watchdog-sla] ingestion started")
	_, _ = pipeline.New("trigger-watchdog-sla", src).
		Then(then).
		To(sink).
		Run(context.Background())
	fmt.Fprintln(os.Stdout, "[trigger-watchdog-sla] ingestion finished")

	// Requirement: example exits 0 even when alerting.
	os.Exit(0)
}
