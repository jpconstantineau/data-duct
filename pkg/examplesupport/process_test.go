package examplesupport

import (
	"context"
	"testing"
	"time"
)

func TestProcessStep_HandleBuildsRequestFromTriggerEvent(t *testing.T) {
	fixed := time.Date(2025, 12, 26, 12, 0, 0, 0, time.UTC)

	p := ProcessStep{
		Now: func() time.Time { return fixed },
		Process: func(ctx context.Context, req IngestionRequest) (IngestionResult, error) {
			if req.Reason != "webhook" {
				t.Fatalf("expected Reason=webhook, got %q", req.Reason)
			}
			if req.InputRef != "local://test" {
				t.Fatalf("expected InputRef=local://test, got %q", req.InputRef)
			}
			if !req.RequestedAt.Equal(fixed) {
				t.Fatalf("expected RequestedAt=%v, got %v", fixed, req.RequestedAt)
			}
			return IngestionResult{StartedAt: fixed, FinishedAt: fixed, Records: 1, Outcome: "succeeded"}, nil
		},
	}

	ev := TriggerEvent{Kind: "webhook", Occurred: fixed, SourceRef: "local://test"}
	_, err := p.Handle(context.Background(), ev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
