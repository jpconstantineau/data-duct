package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jpconstantineau/data-duct/pkg/examplesupport"
	"github.com/jpconstantineau/data-duct/pkg/pipeline"
)

func main() {
	var input = flag.String("input", "", "path to exported warehouse data file (CSV/JSON)")
	flag.Parse()

	if *input == "" {
		fmt.Fprintln(os.Stderr, "missing -input")
		os.Exit(2)
	}

	absInput, err := filepath.Abs(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid -input %q: %v\n", *input, err)
		os.Exit(2)
	}

	process := examplesupport.ProcessStep{
		Process: func(ctx context.Context, req examplesupport.IngestionRequest) (examplesupport.IngestionResult, error) {
			start := time.Now()

			f, err := os.Open(absInput)
			if err != nil {
				return examplesupport.IngestionResult{}, err
			}
			defer f.Close()

			records, err := examplesupport.CountCSVRecords(f, true)
			if err != nil {
				return examplesupport.IngestionResult{}, err
			}

			end := time.Now()
			return examplesupport.IngestionResult{StartedAt: start, FinishedAt: end, Records: records, Outcome: "succeeded"}, nil
		},
	}

	ev := examplesupport.TriggerEvent{Kind: "manual", Occurred: time.Now(), SourceRef: absInput}

	src := func(ctx context.Context) (<-chan examplesupport.TriggerEvent, error) {
		ch := make(chan examplesupport.TriggerEvent, 1)
		ch <- ev
		close(ch)
		return ch, nil
	}

	sink := func(ctx context.Context, res examplesupport.IngestionResult) error {
		examplesupport.PrintOutcome(os.Stdout, "source-warehouse-export", res.Outcome, fmt.Sprintf("rows=%d input=%s", res.Records, absInput))
		return nil
	}

	fmt.Fprintln(os.Stdout, "[source-warehouse-export] trigger fired")
	fmt.Fprintln(os.Stdout, "[source-warehouse-export] ingestion started")
	_, err = pipeline.New("source-warehouse-export", src).
		Then(process.Handle).
		To(sink).
		Run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "[source-warehouse-export] error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, "[source-warehouse-export] ingestion finished")
}
