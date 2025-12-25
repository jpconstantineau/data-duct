package pipelineinternal

import (
	"context"
	"errors"
	"sync"
	"time"
)

type RunState int

const (
	StateSucceeded RunState = iota
	StateCancelled
	StateFailed
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type nopLogger struct{}

func (nopLogger) Debug(string, ...any) {}
func (nopLogger) Info(string, ...any)  {}
func (nopLogger) Warn(string, ...any)  {}
func (nopLogger) Error(string, ...any) {}

type Config struct {
	DefaultBuffer int
	Logger        Logger
}

type StageKind int

const (
	StageSingle StageKind = iota
	StageBatch
)

type StageConfig struct {
	Buffer      int
	Concurrency int
	Name        string
}

type BatchPolicy struct {
	Size    int
	MaxWait time.Duration
}

type SingleHandler func(ctx context.Context, input any) (any, error)

type BatchHandler func(ctx context.Context, inputs []any) ([]any, error)

type Stage struct {
	Kind        StageKind
	Config      StageConfig
	Single      SingleHandler
	Batch       BatchHandler
	BatchPolicy BatchPolicy
}

type Source func(ctx context.Context) (<-chan any, error)

type Sink func(ctx context.Context, input any) error

type feed struct {
	RootCtx      context.Context
	PipelineName string
	Data         any
}

// Run executes the pipeline and blocks until all internal goroutines exit.
func Run(rootCtx context.Context, pipelineName string, source Source, stages []Stage, sink Sink, cfg Config) (RunState, error) {
	if rootCtx == nil {
		rootCtx = context.Background()
	}
	if source == nil || sink == nil {
		return StateFailed, ErrInvalidConfig
	}
	logger := cfg.Logger
	if logger == nil {
		logger = nopLogger{}
	}

	sourceCtx, cancelSource := context.WithCancelCause(rootCtx)
	defer cancelSource(nil)

	policy := &errorPolicy{}

	// Start source.
	srcCh, err := source(sourceCtx)
	if err != nil {
		return StateFailed, err
	}

	// Pump source into first stage as feed.
	in0 := make(chan feed, max(0, cfg.DefaultBuffer))
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(in0)
		sourcePump(rootCtx, sourceCtx, srcCh, in0, pipelineName)
	}()

	// Wire stages.
	current := (<-chan feed)(in0)
	for i := range stages {
		st := stages[i]
		if st.Config.Concurrency < 1 {
			st.Config.Concurrency = 1
		}
		buf := st.Config.Buffer
		if buf < 0 {
			buf = cfg.DefaultBuffer
		}

		out := make(chan feed, max(0, buf))

		switch st.Kind {
		case StageBatch:
			if st.Batch == nil {
				return StateFailed, ErrInvalidConfig
			}
			wg.Add(1)
			go func(in <-chan feed, out chan<- feed, st Stage) {
				defer wg.Done()
				workerBatch(rootCtx, in, out, safeBatch(st.Config.Name, st.Batch, policy), st.BatchPolicy, logger)
			}(current, out, st)
		default:
			if st.Single == nil {
				return StateFailed, ErrInvalidConfig
			}
			wg.Add(1)
			go func(in <-chan feed, out chan<- feed, st Stage) {
				defer wg.Done()
				workerSingle(rootCtx, in, out, safeSingle(st.Config.Name, st.Single, policy), st.Config.Concurrency, logger)
			}(current, out, st)
		}

		current = out
	}

	// Consume sink in the current goroutine to simplify shutdown.
	sinkConsume(rootCtx, current, sink, policy, logger)

	// Stop feeding the source promptly once sink is done.
	if cause := policy.get(); cause != nil {
		cancelSource(cause)
	}

	wg.Wait()

	// Cancellation wins.
	if err := rootCtx.Err(); err != nil {
		return StateCancelled, err
	}

	if cause := policy.get(); cause != nil {
		// If the source context was canceled with a cause, normalize to the cause.
		if errors.Is(cause, context.Canceled) || errors.Is(cause, context.DeadlineExceeded) {
			return StateCancelled, cause
		}
		return StateFailed, cause
	}

	return StateSucceeded, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
