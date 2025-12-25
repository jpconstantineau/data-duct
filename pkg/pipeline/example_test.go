package pipeline

import (
	"context"
	"fmt"
)

func Example() {
	src := func(ctx context.Context) (<-chan int, error) {
		ch := make(chan int, 3)
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)
		return ch, nil
	}

	var out []string
	sink := func(ctx context.Context, s string) error {
		out = append(out, s)
		return nil
	}

	res, _ := New("example", src).
		Then(func(ctx context.Context, n int) (string, error) { return fmt.Sprintf("%d", n), nil }).
		To(sink).
		Run(context.Background())

	fmt.Println(res.State())
	fmt.Println(out)
	// Output:
	// succeeded
	// [1 2 3]
}
