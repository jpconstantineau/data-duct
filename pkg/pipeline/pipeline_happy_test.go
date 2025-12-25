package pipeline

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestPipelineHappyPath_SourceThenTo(t *testing.T) {
	t.Parallel()

	src := func(ctx context.Context) (<-chan int, error) {
		ch := make(chan int, 3)
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)
		return ch, nil
	}

	var got []string
	sink := func(ctx context.Context, s string) error {
		got = append(got, s)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := New("happy", src).
		Then(func(ctx context.Context, n int) (string, error) {
			return "n=" + itoa(n*10), nil
		}).
		To(sink).
		Run(ctx)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res.State() != StateSucceeded {
		t.Fatalf("expected succeeded, got %s", res.State())
	}

	want := []string{"n=10", "n=20", "n=30"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func itoa(n int) string {
	// tiny local helper to keep tests stdlib-only and avoid fmt formatting differences
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	var b [32]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + (n % 10))
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
