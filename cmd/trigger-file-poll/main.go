package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/jpconstantineau/data-duct/pkg/examplesupport"
	"github.com/jpconstantineau/data-duct/pkg/pipeline"
)

func main() {
	var (
		watchDir = flag.String("watch", ".", "directory to watch")
		pattern  = flag.String("pattern", "", "optional pattern matched against base filename (filepath.Match)")
		poll     = flag.Duration("poll", 1*time.Second, "poll interval")
	)
	flag.Parse()

	absWatch, err := filepath.Abs(*watchDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid -watch %q: %v\n", *watchDir, err)
		os.Exit(2)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()

	coalescer := examplesupport.NewCoalescer()
	poller := examplesupport.NewFilePoller(examplesupport.FilePollConfig{WatchDir: absWatch, Pattern: *pattern})

	process := examplesupport.ProcessStep{
		Process: func(ctx context.Context, req examplesupport.IngestionRequest) (examplesupport.IngestionResult, error) {
			start := time.Now()
			// Simulate work.
			time.Sleep(200 * time.Millisecond)
			end := time.Now()
			return examplesupport.IngestionResult{StartedAt: start, FinishedAt: end, Records: 1, Outcome: "succeeded"}, nil
		},
	}

	var (
		mu        sync.Mutex
		lastEvent examplesupport.TriggerEvent
		hasEvent  bool
	)

	runnerDone := make(chan error, 1)
	go func() {
		runnerDone <- coalescer.Run(ctx, func(ctx context.Context) error {
			mu.Lock()
			ev := lastEvent
			have := hasEvent
			hasEvent = false
			mu.Unlock()
			if !have {
				ev = examplesupport.TriggerEvent{Kind: "file_poll", Occurred: time.Now(), SourceRef: "(unknown)"}
			}

			fmt.Fprintln(os.Stdout, "[trigger-file-poll] ingestion started")

			src := func(ctx context.Context) (<-chan examplesupport.TriggerEvent, error) {
				ch := make(chan examplesupport.TriggerEvent, 1)
				ch <- ev
				close(ch)
				return ch, nil
			}

			sink := func(ctx context.Context, res examplesupport.IngestionResult) error {
				examplesupport.PrintOutcome(os.Stdout, "trigger-file-poll", res.Outcome, fmt.Sprintf("records=%d", res.Records))
				fmt.Fprintln(os.Stdout, "[trigger-file-poll] ingestion finished")
				return nil
			}

			_, err := pipeline.New("trigger-file-poll", src).
				Then(process.Handle).
				To(sink).
				Run(ctx)
			return err
		})
	}()

	fmt.Fprintf(os.Stdout, "[trigger-file-poll] watching dir=%s pattern=%q poll=%s\n", absWatch, *pattern, poll.String())

	ticker := time.NewTicker(*poll)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			cancel()
			_ = <-runnerDone
			return
		case <-ticker.C:
			ev, ok, err := poller.Scan()
			if err != nil {
				fmt.Fprintf(os.Stderr, "[trigger-file-poll] scan error: %v\n", err)
				continue
			}
			if !ok {
				continue
			}

			fmt.Fprintf(os.Stdout, "[trigger-file-poll] trigger fired file=%s\n", ev.SourceRef)
			mu.Lock()
			lastEvent = ev
			hasEvent = true
			mu.Unlock()

			coalescer.Request()
		}
	}
}
