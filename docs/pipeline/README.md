# Graceful Context Pipeline

This package provides a small, stdlib-only Go library for building and running a typed, channel-based pipeline that respects context cancellation and shuts down gracefully.

## Concepts

- **Source**: `func(ctx context.Context) (<-chan T, error)`
- **Processor**: `Handler[In, Out]` (`func(ctx context.Context, input In) (Out, error)`)
- **Batch processor**: `BatchHandler[In, Out]` (`func(ctx context.Context, inputs []In) ([]Out, error)`)
- **Sink**: `EndHandler[T]` (`func(ctx context.Context, input T) error`)

## Quick example

See the runnable example in `cmd/graceful-context-pipeline-example`.

## Cancellation & errors

- Root context cancellation stops the pipeline and returns a `Cancelled` result.
- Processor/sink errors stop acceptance of new inputs and return a `Failed` result.
- Panics in user handlers are recovered and returned as errors.

## Batching

Use `ThenBatch` with a `BatchPolicy` to group items into deterministic batches. The current implementation supports:

- Fixed `Size`
- Optional `MaxWait` to flush early

## Commands

```powershell
go test ./...
go run ./cmd/graceful-context-pipeline-example

# trigger & source examples (see specs/001-trigger-source-examples/quickstart.md)
go run ./cmd/trigger-interval
go run ./cmd/trigger-daily
go run ./cmd/trigger-webhook
go run ./cmd/trigger-file-poll
go run ./cmd/trigger-manual-run-now
go run ./cmd/trigger-watchdog-sla

go run ./cmd/source-files
go run ./cmd/source-database-export
go run ./cmd/source-warehouse-export
go run ./cmd/source-rest-uselessfacts
go run ./cmd/source-observability

# lint (optional)
./scripts/ci/lint.ps1

# vulnerability scan (optional)
./scripts/ci/vuln.ps1
```

