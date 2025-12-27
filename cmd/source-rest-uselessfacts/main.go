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
	var url = flag.String("url", "https://uselessfacts.jsph.pl/api/v2/facts/random?language=en", "uselessfacts API URL")
	flag.Parse()

	process := examplesupport.ProcessStep{
		Process: func(ctx context.Context, req examplesupport.IngestionRequest) (examplesupport.IngestionResult, error) {
			start := time.Now()

			fact, err := examplesupport.FetchUselessFactText(ctx, *url)
			if err != nil {
				return examplesupport.IngestionResult{}, err
			}

			end := time.Now()
			// Records=1 because we ingested one fact.
			_ = fact
			return examplesupport.IngestionResult{StartedAt: start, FinishedAt: end, Records: 1, Outcome: "succeeded"}, nil
		},
	}

	ev := examplesupport.TriggerEvent{Kind: "manual", Occurred: time.Now(), SourceRef: *url}

	src := func(ctx context.Context) (<-chan examplesupport.TriggerEvent, error) {
		ch := make(chan examplesupport.TriggerEvent, 1)
		ch <- ev
		close(ch)
		return ch, nil
	}

	sink := func(ctx context.Context, res examplesupport.IngestionResult) error {
		examplesupport.PrintOutcome(os.Stdout, "source-rest-uselessfacts", res.Outcome, fmt.Sprintf("records=%d url=%s", res.Records, *url))
		return nil
	}

	fmt.Fprintln(os.Stdout, "[source-rest-uselessfacts] trigger fired")
	fmt.Fprintln(os.Stdout, "[source-rest-uselessfacts] ingestion started")
	_, err := pipeline.New("source-rest-uselessfacts", src).
		Then(process.Handle).
		To(sink).
		Run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "[source-rest-uselessfacts] error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, "[source-rest-uselessfacts] ingestion finished")
}
