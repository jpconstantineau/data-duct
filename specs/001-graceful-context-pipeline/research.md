# Research: Graceful Context Pipeline

**Date**: 2025-12-24

This document captures the key technical decisions for implementing the feature, with rationale and alternatives.

## Decision: Source contract returns a read-only channel

- **Decision**: Source is `func(ctx context.Context) (<-chan T, error)`.
- **Rationale**: Minimizes user boilerplate, matches the channel-based pipeline concept directly, and makes completion explicit via channel close.
- **Alternatives considered**:
  - `func(ctx context.Context, out chan<- T) error`: more flexible but requires extra ownership rules and more boilerplate.
  - `func(ctx context.Context) ([]T, error)`: not streaming-friendly.

## Decision: In-flight handling on error/cancellation

- **Decision**: Stop accepting new inputs promptly, cancel upstream work, and best-effort finish items already buffered/in-flight inside the pipeline, bounded by the root context.
- **Rationale**: Predictable shutdown without forcing a potentially unbounded drain; minimizes user surprise and avoids hangs.
- **Alternatives considered**:
  - Fail-fast drop everything: simplest but can lose buffered work.
  - Full drain after error: can hang and violates graceful shutdown expectations.

## Decision: User handlers operate on raw payload `T`

- **Decision**: User processor and sink handlers take/return raw `T` (or `[]T` for batch), not `Feed[T]`.
- **Rationale**: Keeps user code minimal and focused on business logic. Metadata stays internal.
- **Alternatives considered**:
  - Expose `Feed[T]` in user handler signatures: increases boilerplate and couples user code to pipeline transport.

## Decision: Default concurrency

- **Decision**: Default stage concurrency to 1 worker per stage.
- **Rationale**: Safe baseline; avoids surprising parallelism. Advanced users opt-in.
- **Alternatives considered**:
  - Default to CPU count: can overwhelm downstreams and makes ordering/side-effects harder.
  - Require explicit configuration: forces boilerplate.

## Decision: Structured logging

- **Decision**: If logging is needed, use Go standard library `log/slog`.
- **Rationale**: Structured logging without runtime dependencies; aligns with constitution core constraints.
- **Alternatives considered**:
  - Third-party loggers: disallowed for `core` runtime dependencies.

## Decision: Package layout (project-layout)

- **Decision**: Public API in `pkg/pipeline`; internal orchestration in `internal/pipelineinternal`; example binary in `cmd/graceful-context-pipeline-example`.
- **Rationale**: Matches golang-standards/project-layout and keeps internal machinery hidden.
- **Alternatives considered**:
  - Single `internal/` only: makes public API less discoverable.
  - `pkg/core/`: unclear naming; this feature is specifically pipeline orchestration.

## Decision: Fluent builder API with build-time signature validation

- **Decision**: Provide a fluent builder style that reads naturally left-to-right: `New(...).Then(...).To(...).Run(...)`.
- **Rationale**: This is the most readable and ergonomic shape for typical Go users and matches how pipeline stages are conceptualized.
- **Constraint encountered**: Go does not support generic methods on generic types, which prevents an idealized fully generic `(*Pipeline[T]).Then[U](...) *Pipeline[U]` fluent API.
- **Implementation approach chosen**: Accept ordinary typed handler functions and validate their signatures and stage compatibility when stages are added, failing fast on invalid wiring.
- **Alternatives considered**:
  - Function-based generic API (`Then(p, handler)`) to keep compile-time type chaining: correct but less readable and explicitly disallowed by the desired example style.
  - Interface-based untyped handlers (`func(context.Context, any) (any, error)`): simplest but pushes type casting to users and makes misuse easier.
