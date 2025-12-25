package pipeline

import "log/slog"

type Option func(*pipelineOptions)

type StageOption func(*stageOptions)

type pipelineOptions struct {
	buffer int
	logger *slog.Logger
}

type stageOptions struct {
	buffer      int
	concurrency int
	name        string
}

func defaultPipelineOptions() pipelineOptions {
	return pipelineOptions{buffer: 0, logger: nil}
}

func defaultStageOptions() stageOptions {
	return stageOptions{buffer: 0, concurrency: 1, name: ""}
}

// WithBuffer sets the default inter-stage buffer size.
func WithBuffer(n int) Option {
	return func(o *pipelineOptions) {
		if n < 0 {
			n = 0
		}
		o.buffer = n
	}
}

// WithLogger enables optional structured logging.
func WithLogger(logger *slog.Logger) Option {
	return func(o *pipelineOptions) {
		o.logger = logger
	}
}

// WithStageBuffer sets the buffer size between this stage and the next.
func WithStageBuffer(n int) StageOption {
	return func(o *stageOptions) {
		if n < 0 {
			n = 0
		}
		o.buffer = n
	}
}

// WithStageConcurrency sets the number of workers for a stage.
func WithStageConcurrency(n int) StageOption {
	return func(o *stageOptions) {
		if n < 1 {
			n = 1
		}
		o.concurrency = n
	}
}

// WithStageName labels a stage (primarily for logging).
func WithStageName(name string) StageOption {
	return func(o *stageOptions) {
		o.name = name
	}
}
