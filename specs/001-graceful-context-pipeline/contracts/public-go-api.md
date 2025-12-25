# Contract: Public Go API (Draft)

**Date**: 2025-12-24

This contract describes the intended public API surface for the feature. It is not an HTTP/REST API.

## Handler Types

- `Handler[In any, Out any]`: `func(ctx context.Context, input In) (Out, error)`
- `BatchHandler[In any, Out any]`: `func(ctx context.Context, inputs []In) ([]Out, error)`
- `StartHandler[T any]`: `func(ctx context.Context) (T, error)`
- `EndHandler[T any]`: `func(ctx context.Context, input T) error`

Note: These generic handler type aliases describe the *intended shapes* of user-provided functions. The fluent builder accepts ordinary typed functions directly.

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

- Fluent, readable stage chaining (`New(...).Then(...).To(...)`)
- Strong stage-to-stage compatibility checks when building the pipeline (fail fast on invalid wiring)
- Minimal boilerplate for users
- Defaults: concurrency=1 per stage

Implemented shape:

- `pipeline.New[T](name string, source func(context.Context) (<-chan T, error), opts ...Option) *Pipeline`
- `(*Pipeline).Then(handler func(context.Context, In) (Out, error), opts ...StageOption) *Pipeline`
- `(*Pipeline).ThenBatch(handler func(context.Context, []In) ([]Out, error), batch BatchPolicy, opts ...StageOption) *Pipeline`
- `(*Pipeline).To(sink func(context.Context, In) error, opts ...StageOption) *Runnable`
- `(*Runnable).Run(ctx context.Context) (Result, error)`

Validation:

- The library validates each handler’s function signature when added.
- The library validates that a stage’s input type matches the previous stage’s output type.
- Invalid wiring fails fast (panic) during pipeline construction.

## Cancellation / Error Semantics

- Default error policy: stop-on-first-error.
- On error/cancellation: stop accepting new inputs, cancel upstream work, and best-effort finish items already buffered/in-flight inside the pipeline, bounded by the root context.

## Non-Goals (for MVP)

- Exactly-once delivery guarantees
- Persistent checkpointing
- Automatic retries/backoff
- Ordering guarantees by default
