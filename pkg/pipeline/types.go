package pipeline

import "context"

// Feed is the conceptual transport object that associates pipeline metadata
// with a payload. User handlers operate on raw payload types, not Feed.
type Feed[T any] struct {
	RootCtx      context.Context
	PipelineName string
	Data         T
}

type SourceFunc[T any] func(ctx context.Context) (<-chan T, error)

type Handler[In any, Out any] func(ctx context.Context, input In) (Out, error)

type BatchHandler[In any, Out any] func(ctx context.Context, inputs []In) ([]Out, error)

type EndHandler[T any] func(ctx context.Context, input T) error
