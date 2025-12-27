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
	var runNow = flag.Bool("run-now", false, "trigger an ingestion run immediately")
	flag.Parse()

	ev, ok := examplesupport.ManualRunNow(*runNow, time.Now)
	if !ok {
		fmt.Fprintln(os.Stdout, "[trigger-manual-run-now] pass -run-now to trigger ingestion")
		return
	}

	process := examplesupport.ProcessStep{
		Process: func(ctx context.Context, req examplesupport.IngestionRequest) (examplesupport.IngestionResult, error) {
			start := time.Now()
			time.Sleep(100 * time.Millisecond)
			end := time.Now()
			return examplesupport.IngestionResult{StartedAt: start, FinishedAt: end, Records: 1, Outcome: "succeeded"}, nil
		},
	}

	src := func(ctx context.Context) (<-chan examplesupport.TriggerEvent, error) {
		ch := make(chan examplesupport.TriggerEvent, 1)
		ch <- ev
		close(ch)
		return ch, nil
	}

	sink := func(ctx context.Context, res examplesupport.IngestionResult) error {
		examplesupport.PrintOutcome(os.Stdout, "trigger-manual-run-now", res.Outcome, fmt.Sprintf("records=%d", res.Records))
		return nil
	}

	fmt.Fprintln(os.Stdout, "[trigger-manual-run-now] trigger fired")
	fmt.Fprintln(os.Stdout, "[trigger-manual-run-now] ingestion started")
	_, err := pipeline.New("trigger-manual-run-now", src).
		Then(process.Handle).
		To(sink).
		Run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "[trigger-manual-run-now] error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, "[trigger-manual-run-now] ingestion finished")
}
