package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jpconstantineau/data-duct/pkg/examplesupport"
	"github.com/jpconstantineau/data-duct/pkg/pipeline"
)

func main() {
	var (
		logs      = flag.String("logs", "", "path to log file")
		traces    = flag.String("traces", "", "path to trace export JSON file")
		metricsUR = flag.String("metrics-url", "", "optional URL returning metrics JSON")
	)
	flag.Parse()

	if *logs == "" || *traces == "" {
		fmt.Fprintln(os.Stderr, "missing required -logs and/or -traces")
		os.Exit(2)
	}

	absLogs, err := filepath.Abs(*logs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid -logs %q: %v\n", *logs, err)
		os.Exit(2)
	}
	absTraces, err := filepath.Abs(*traces)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid -traces %q: %v\n", *traces, err)
		os.Exit(2)
	}

	process := examplesupport.ProcessStep{
		Process: func(ctx context.Context, req examplesupport.IngestionRequest) (examplesupport.IngestionResult, error) {
			start := time.Now()

			lf, err := os.Open(absLogs)
			if err != nil {
				return examplesupport.IngestionResult{}, err
			}
			defer lf.Close()
			logLines, err := examplesupport.CountLogLines(lf)
			if err != nil {
				return examplesupport.IngestionResult{}, err
			}

			tf, err := os.Open(absTraces)
			if err != nil {
				return examplesupport.IngestionResult{}, err
			}
			defer tf.Close()
			traceSpans, err := examplesupport.CountJSONArrayItems(tf)
			if err != nil {
				return examplesupport.IngestionResult{}, err
			}

			metricsCount := 0
			if *metricsUR != "" {
				mc, err := fetchJSONArrayCount(ctx, *metricsUR)
				if err != nil {
					return examplesupport.IngestionResult{}, err
				}
				metricsCount = mc
			}

			end := time.Now()
			total := logLines + traceSpans + metricsCount
			return examplesupport.IngestionResult{StartedAt: start, FinishedAt: end, Records: total, Outcome: "succeeded"}, nil
		},
	}

	ev := examplesupport.TriggerEvent{Kind: "manual", Occurred: time.Now(), SourceRef: "observability://local"}

	src := func(ctx context.Context) (<-chan examplesupport.TriggerEvent, error) {
		ch := make(chan examplesupport.TriggerEvent, 1)
		ch <- ev
		close(ch)
		return ch, nil
	}

	sink := func(ctx context.Context, res examplesupport.IngestionResult) error {
		summary := fmt.Sprintf("logs=%d traces=%s metrics_url=%q", res.Records, absTraces, *metricsUR)
		examplesupport.PrintOutcome(os.Stdout, "source-observability", res.Outcome, summary)
		return nil
	}

	fmt.Fprintln(os.Stdout, "[source-observability] trigger fired")
	fmt.Fprintln(os.Stdout, "[source-observability] ingestion started")
	_, err = pipeline.New("source-observability", src).
		Then(process.Handle).
		To(sink).
		Run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "[source-observability] error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, "[source-observability] ingestion finished")
}

func fetchJSONArrayCount(ctx context.Context, url string) (int, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return 0, fmt.Errorf("metrics url returned status %d", res.StatusCode)
	}

	lr := &io.LimitedReader{R: res.Body, N: 1024 * 1024}
	return examplesupport.CountJSONArrayItems(lr)
}
