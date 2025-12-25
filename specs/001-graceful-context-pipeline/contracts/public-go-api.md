# Contract: Public Go API (Draft)

**Date**: 2025-12-24

This contract describes the intended public API surface for the feature. It is not an HTTP/REST API.

## Handler Types

- `Handler[In any, Out any]`: `func(ctx context.Context, input In) (Out, error)`
- `BatchHandler[In any, Out any]`: `func(ctx context.Context, inputs []In) ([]Out, error)`
- `StartHandler[T any]`: `func(ctx context.Context) (T, error)`
- `EndHandler[T any]`: `func(ctx context.Context, input T) error`

## Data Transport

- Internal transport uses `Feed[T]`:
  - `RootCtx context.Context`
  - `PipelineName string`
  - `Data T`

- User handlers operate on raw `T` (not `Feed[T]`).

## Source Contract

- Source MUST be `func(ctx context.Context) (<-chan T, error)`.
- Source MUST close the returned channel when complete.

## Pipeline Construction API (Intent)

Goals:

- Compile-time type safety across stages (`T0 -> T1 -> ... -> TN`)
- Minimal boilerplate for users
- Defaults: concurrency=1 per stage

Candidate shape (subject to refinement during implementation):

- `pipeline.New[T](name string, source func(context.Context) (<-chan T, error), opts ...Option) *Pipeline[T]`
- `(*Pipeline[In]).Then[Out](handler Handler[In, Out], opts ...StageOption) *Pipeline[Out]`
- `(*Pipeline[In]).ThenBatch[Out](handler BatchHandler[In, Out], batch BatchPolicy, opts ...StageOption) *Pipeline[Out]`
- `(*Pipeline[T]).To(sink EndHandler[T], opts ...StageOption) *Runnable`
- `(*Runnable).Run(ctx context.Context) (Result, error)`

## Cancellation / Error Semantics

- Default error policy: stop-on-first-error.
- On error/cancellation: stop accepting new inputs, cancel upstream work, and best-effort finish items already buffered/in-flight inside the pipeline, bounded by the root context.

## Non-Goals (for MVP)

- Exactly-once delivery guarantees
- Persistent checkpointing
- Automatic retries/backoff
- Ordering guarantees by default
