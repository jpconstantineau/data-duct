package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jpconstantineau/data-duct/pkg/examplesupport"
	"github.com/jpconstantineau/data-duct/pkg/pipeline"
)

func main() {
	var (
		listenAddr = flag.String("listen", "127.0.0.1:8080", "address to listen on")
		path       = flag.String("path", "/trigger", "webhook path")
	)
	flag.Parse()

	if *path == "" || (*path)[0] != '/' {
		fmt.Fprintf(os.Stderr, "invalid -path %q (must start with /)\n", *path)
		os.Exit(2)
	}

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

	var (
		mu        sync.Mutex
		lastEvent examplesupport.TriggerEvent
		hasEvent  bool
	)

	handler := &examplesupport.WebhookHandler{
		MaxBodyBytes: 64 * 1024,
		OnEvent: func(ctx context.Context, ev examplesupport.TriggerEvent, req examplesupport.WebhookRequest) error {
			mu.Lock()
			lastEvent = ev
			hasEvent = true
			mu.Unlock()

			fmt.Fprintf(os.Stdout, "[trigger-webhook] trigger fired source_category=%q note=%q\n", req.SourceCategory, req.Note)
			coalescer.Request()
			return nil
		},
	}

	runnerDone := make(chan error, 1)
	go func() {
		runnerDone <- coalescer.Run(ctx, func(ctx context.Context) error {
			mu.Lock()
			ev := lastEvent
			have := hasEvent
			hasEvent = false
			mu.Unlock()
			if !have {
				ev = examplesupport.TriggerEvent{Kind: "webhook", Occurred: time.Now(), SourceRef: "webhook://local"}
			}

			fmt.Fprintln(os.Stdout, "[trigger-webhook] ingestion started")

			src := func(ctx context.Context) (<-chan examplesupport.TriggerEvent, error) {
				ch := make(chan examplesupport.TriggerEvent, 1)
				ch <- ev
				close(ch)
				return ch, nil
			}

			sink := func(ctx context.Context, res examplesupport.IngestionResult) error {
				examplesupport.PrintOutcome(os.Stdout, "trigger-webhook", res.Outcome, fmt.Sprintf("records=%d", res.Records))
				fmt.Fprintln(os.Stdout, "[trigger-webhook] ingestion finished")
				return nil
			}

			_, err := pipeline.New("trigger-webhook", src).
				Then(process.Handle).
				To(sink).
				Run(ctx)
			return err
		})
	}()

	mux := http.NewServeMux()
	mux.Handle(*path, handler)

	srv := &http.Server{
		Addr:              *listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Fprintf(os.Stdout, "[trigger-webhook] listening addr=%s path=%s\n", *listenAddr, *path)

	srvErr := make(chan error, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			srvErr <- err
			return
		}
		srvErr <- nil
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sig:
		fmt.Fprintln(os.Stdout, "[trigger-webhook] shutting down")
	case err := <-srvErr:
		if err != nil {
			fmt.Fprintf(os.Stderr, "[trigger-webhook] server error: %v\n", err)
			os.Exit(1)
		}
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 3*time.Second)
	_ = srv.Shutdown(shutdownCtx)
	cancelShutdown()

	cancel()
	_ = <-runnerDone
}
