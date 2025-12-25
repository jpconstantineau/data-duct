package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jpconstantineau/data-duct/pkg/pipeline"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	src := func(ctx context.Context) (<-chan int, error) {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for i := 1; i <= 5; i++ {
				select {
				case <-ctx.Done():
					return
				case ch <- i:
				}
			}
		}()
		return ch, nil
	}

	var out []string
	sink := func(ctx context.Context, s string) error {
		out = append(out, s)
		return nil
	}

	var thn = func(ctx context.Context, n int) (string, error) {
		return fmt.Sprintf("value=%d", n*2), nil
	}

	res, err := pipeline.New("example", src).
		Then(thn).
		To(sink).
		Run(ctx)

	fmt.Printf("result=%s err=%v out=%v\n", res.State(), err, out)
}
