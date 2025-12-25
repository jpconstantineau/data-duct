package pipeline

import (
	"context"

	"github.com/jpconstantineau/data-duct/internal/pipelineinternal"
)

type Runnable struct {
	def *definition
}

func (r *Runnable) Run(ctx context.Context) (Result, error) {
	if r == nil || r.def == nil {
		return Failed{Cause: pipelineinternal.ErrInvalidConfig}, pipelineinternal.ErrInvalidConfig
	}
	if r.def.source == nil || r.def.sink == nil {
		return Failed{Cause: pipelineinternal.ErrInvalidConfig}, pipelineinternal.ErrInvalidConfig
	}

	state, cause := pipelineinternal.Run(
		ctx,
		r.def.name,
		r.def.source,
		toInternalStages(r.def.stages),
		r.def.sink,
		pipelineinternal.Config{DefaultBuffer: r.def.buffer, Logger: r.def.logger},
	)

	switch state {
	case pipelineinternal.StateSucceeded:
		return Succeeded{}, nil
	case pipelineinternal.StateCancelled:
		return Cancelled{Cause: cause}, cause
	default:
		return Failed{Cause: cause}, cause
	}
}

func toInternalStages(stages []stageDef) []pipelineinternal.Stage {
	out := make([]pipelineinternal.Stage, 0, len(stages))
	for _, s := range stages {
		cfg := pipelineinternal.StageConfig{Buffer: s.buffer, Concurrency: s.concurrency, Name: s.name}
		switch s.kind {
		case stageBatch:
			out = append(out, pipelineinternal.Stage{Kind: pipelineinternal.StageBatch, Batch: s.batch, BatchPolicy: s.batchPolicy, Config: cfg})
		default:
			out = append(out, pipelineinternal.Stage{Kind: pipelineinternal.StageSingle, Single: s.single, Config: cfg})
		}
	}
	return out
}
