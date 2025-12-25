# Data Model: Graceful Context Pipeline

**Date**: 2025-12-24

This feature is primarily a concurrency/orchestration library. The "data model" is comprised of conceptual entities and their relationships, expressed as API-level types.

## Entities

### Feed[T]

Represents an item flowing through the pipeline with associated metadata.

- **Fields**:
  - `RootCtx`: context used as the root cancellation/timeout signal for pipeline execution
  - `PipelineName`: string name of the pipeline
  - `Data`: payload item `T`

- **Validation rules**:
  - `RootCtx` MUST be non-nil at runtime (the library will treat nil as `context.Background()` or reject configuration; final choice is specified in contracts).
  - `PipelineName` SHOULD be non-empty (empty allowed but discouraged).

### Pipeline

A configured pipeline definition.

- **Fields (conceptual)**:
  - `Name`
  - `Source` (returns `<-chan T`)
  - `Stages` (processors)
  - `Sink`
  - `Options` (buffer sizes, concurrency per stage, batching policy)

- **Relationships**:
  - Pipeline produces Feeds from Source items
  - Stages transform `T → U` (possibly via batching)
  - Sink consumes the final payload type

### Stage

One step in the pipeline.

- **Types**:
  - Single-item stage: `Handler[In, Out]`
  - Batch stage: `BatchHandler[In, Out]` with batch policy

- **State transitions**:
  - `Created` → `Running` → `Completed | Cancelled | Failed`

### Execution Result

A summary of a pipeline run.

- **Fields (conceptual)**:
  - `State`: {Succeeded, Cancelled, Failed}
  - `Err`: optional error
  - Counters: items in/out per stage (optional; include only if it doesn’t add complexity)

## Notes

- The public API should keep the above types minimal and avoid over-modeling.
- Any optional observability counters should be opt-in and avoid heavy runtime deps.
