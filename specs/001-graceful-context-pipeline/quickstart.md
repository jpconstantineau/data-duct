# Quickstart: Graceful Context Pipeline

**Date**: 2025-12-24

This quickstart describes the intended developer workflow for the feature.

## Goal

Build and run a minimal pipeline:

- Source returns a `<-chan T`
- One processor transforms `T -> U`
- Sink consumes `U`
- Cancellation stops the pipeline gracefully

## Minimal Usage (Conceptual)

1. Define a source function that returns a read-only channel and closes it when done.
2. Add one or more `Then(...)` processors.
3. Attach a sink via `To(...)`.
4. Run with a root context.

## Minimal Usage (Concrete)

```go
src := func(ctx context.Context) (<-chan int, error) { /* ... */ }
proc := func(ctx context.Context, n int) (string, error) { /* ... */ }
sink := func(ctx context.Context, s string) error { /* ... */ }

runnable := pipeline.New("example", src).
	Then(proc).
	To(sink)

res, err := runnable.Run(ctx)
```

## Expected Commands (once implementation exists)

```powershell
# run unit tests
 go test ./...

# format
 gofmt -w .

# security
 govulncheck ./...

# run example
 go run ./cmd/graceful-context-pipeline-example
```

## Example (once implementation exists)

A runnable example will live under:

- `cmd/graceful-context-pipeline-example`

and demonstrate source → processor → sink with cancellation.

