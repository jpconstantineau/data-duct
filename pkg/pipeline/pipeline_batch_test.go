package pipeline

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestPipelineBatchGroupingBehavior(t *testing.T) {
	t.Parallel()

	src := func(ctx context.Context) (<-chan int, error) {
		ch := make(chan int, 5)
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		close(ch)
		return ch, nil
	}

	var batches [][]int
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := New("batch", src).
		ThenBatch(func(ctx context.Context, inputs []int) ([]int, error) {
			// record the batch
			cp := append([]int(nil), inputs...)
			batches = append(batches, cp)
			// passthrough
			return inputs, nil
		}, BatchPolicy{Size: 2}).
		To(func(ctx context.Context, n int) error { return nil }).
		Run(ctx)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res.State() != StateSucceeded {
		t.Fatalf("expected succeeded, got %s", res.State())
	}

	want := [][]int{{1, 2}, {3, 4}, {5}}
	if !reflect.DeepEqual(batches, want) {
		t.Fatalf("got %v want %v", batches, want)
	}
}
