package examplesupport

import (
	"testing"
	"time"
)

func TestManualRunNow_NotRequested(t *testing.T) {
	if _, ok := ManualRunNow(false, func() time.Time { return time.Unix(1, 0).UTC() }); ok {
		t.Fatalf("expected ok=false")
	}
}

func TestManualRunNow_Requested(t *testing.T) {
	now := time.Unix(123, 0).UTC()
	ev, ok := ManualRunNow(true, func() time.Time { return now })
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if ev.Kind != "manual" {
		t.Fatalf("expected kind manual, got %q", ev.Kind)
	}
	if !ev.Occurred.Equal(now) {
		t.Fatalf("expected occurred %v, got %v", now, ev.Occurred)
	}
	if ev.SourceRef != "manual://run-now" {
		t.Fatalf("expected sourceRef manual://run-now, got %q", ev.SourceRef)
	}
}
