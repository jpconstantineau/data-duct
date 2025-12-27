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
	var root = flag.String("root", ".", "root directory (local folder stand-in for files/object storage)")
	flag.Parse()

	absRoot, err := filepath.Abs(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid -root %q: %v\n", *root, err)
		os.Exit(2)
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot stat root %q: %v\n", absRoot, err)
		os.Exit(2)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "root %q is not a directory\n", absRoot)
		os.Exit(2)
	}

	process := examplesupport.ProcessStep{
		Process: func(ctx context.Context, req examplesupport.IngestionRequest) (examplesupport.IngestionResult, error) {
			start := time.Now()

			entries, err := os.ReadDir(absRoot)
			if err != nil {
				return examplesupport.IngestionResult{}, err
			}
			files := 0
			for _, e := range entries {
				if !e.IsDir() {
					files++
				}
			}

			end := time.Now()
			return examplesupport.IngestionResult{StartedAt: start, FinishedAt: end, Records: files, Outcome: "succeeded"}, nil
		},
	}

	ev := examplesupport.TriggerEvent{Kind: "manual", Occurred: time.Now(), SourceRef: absRoot}

	src := func(ctx context.Context) (<-chan examplesupport.TriggerEvent, error) {
		ch := make(chan examplesupport.TriggerEvent, 1)
		ch <- ev
		close(ch)
		return ch, nil
	}

	sink := func(ctx context.Context, res examplesupport.IngestionResult) error {
		examplesupport.PrintOutcome(os.Stdout, "source-files", res.Outcome, fmt.Sprintf("files=%d root=%s", res.Records, absRoot))
		return nil
	}

	fmt.Fprintln(os.Stdout, "[source-files] trigger fired")
	fmt.Fprintln(os.Stdout, "[source-files] ingestion started")
	_, err = pipeline.New("source-files", src).
		Then(process.Handle).
		To(sink).
		Run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "[source-files] error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, "[source-files] ingestion finished")
}
