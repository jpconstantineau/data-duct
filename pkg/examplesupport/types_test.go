package examplesupport

import (
	"testing"
	"time"
)

func TestTriggerEventValidate(t *testing.T) {
	var ev TriggerEvent
	if err := ev.Validate(); err == nil {
		t.Fatalf("expected Validate() error for zero TriggerEvent")
	}

	ev = TriggerEvent{Kind: "interval", Occurred: time.Now()}
	if err := ev.Validate(); err != nil {
		t.Fatalf("expected Validate() success, got: %v", err)
	}
}

func TestIngestionRequestValidate(t *testing.T) {
	var req IngestionRequest
	if err := req.Validate(); err == nil {
		t.Fatalf("expected Validate() error for zero IngestionRequest")
	}

	req = IngestionRequest{RequestedAt: time.Now(), Reason: "manual"}
	if err := req.Validate(); err != nil {
		t.Fatalf("expected Validate() success, got: %v", err)
	}
}
