package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jpconstantineau/data-duct/pkg/pipeline"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	src := func(ctx context.Context) (<-chan int, error) {
		ch := make(chan int)
		go func() {
			defer close(ch)

			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()

			timer := time.NewTimer(10 * time.Second)
			defer timer.Stop()

			n := 0
			for {
				select {
				case <-ctx.Done():
					return
				case <-timer.C:
					return
				case <-ticker.C:
					n++
					select {
					case <-ctx.Done():
						return
					case ch <- n:
					}
				}
			}
		}()
		return ch, nil
	}

	batchSum := func(ctx context.Context, inputs []int) ([]int, error) {
		sum := 0
		for _, v := range inputs {
			sum += v
		}
		// Emit one output per batch: the sum.
		return []int{sum}, nil
	}

	batchCount := 0
	sink := func(ctx context.Context, batchTotal int) error {
		batchCount++
		fmt.Printf("batch=%d sum=%d\n", batchCount, batchTotal)
		return nil
	}

	runnable := pipeline.New("batch-simple", src).
		ThenBatch(batchSum, pipeline.BatchPolicy{Size: 10}).
		To(sink)

	res, err := runnable.Run(ctx)
	fmt.Printf("result=%s err=%v batches=%d\n", res.State(), err, batchCount)
}
