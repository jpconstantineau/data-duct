package examplesupport

import (
	"context"
	"errors"
	"time"
)

var ErrInvalidProcessStep = errors.New("invalid process step")

// ProcessStep adapts a TriggerEvent into an IngestionRequest and calls Process.
//
// Runnable examples can plug in datasource-specific ingestion logic via Process.
type ProcessStep struct {
	Now     func() time.Time
	Process func(ctx context.Context, req IngestionRequest) (IngestionResult, error)
}

func (p ProcessStep) Handle(ctx context.Context, ev TriggerEvent) (IngestionResult, error) {
	if err := ev.Validate(); err != nil {
		return IngestionResult{}, err
	}
	if p.Process == nil {
		return IngestionResult{}, ErrInvalidProcessStep
	}

	now := p.Now
	if now == nil {
		now = time.Now
	}

	req := IngestionRequest{
		RequestedAt: now(),
		Reason:      ev.Kind,
		InputRef:    ev.SourceRef,
	}
	if err := req.Validate(); err != nil {
		return IngestionResult{}, err
	}

	return p.Process(ctx, req)
}
