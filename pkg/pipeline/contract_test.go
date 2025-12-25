package pipeline

import (
	"context"
	"testing"
)

// This file ensures the fluent API is compile-time type safe.

func compileTimeSource[T any](items []T) SourceFunc[T] {
	return func(ctx context.Context) (<-chan T, error) {
		ch := make(chan T, len(items))
		for _, it := range items {
			ch <- it
		}
		close(ch)
		return ch, nil
	}
}

func TestChainBuild_SignaturesAccepted(t *testing.T) {
	_ = New("ct", compileTimeSource([]int{1, 2})).
		Then(func(ctx context.Context, in int) (string, error) { return "x", nil }).
		Then(func(ctx context.Context, in string) (int, error) { return len(in), nil }).
		To(func(ctx context.Context, in int) error { return nil })
}
