package examplesupport

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidTriggerEvent     = errors.New("invalid trigger event")
	ErrInvalidIngestionRequest = errors.New("invalid ingestion request")
)

// TriggerEvent represents a trigger occurrence (schedule tick, webhook event, file detection, etc.).
//
// Examples normalize all trigger inputs into this shape before constructing an IngestionRequest.
type TriggerEvent struct {
	Kind      string
	Occurred  time.Time
	Payload   []byte
	SourceRef string
}

func (e TriggerEvent) Validate() error {
	if strings.TrimSpace(e.Kind) == "" {
		return ErrInvalidTriggerEvent
	}
	if e.Occurred.IsZero() {
		return ErrInvalidTriggerEvent
	}
	return nil
}

// IngestionRequest is a normalized request passed to the process step.
// It intentionally stays vendor-neutral and dependency-free.
type IngestionRequest struct {
	RequestedAt time.Time
	Reason      string
	InputRef    string
}

func (r IngestionRequest) Validate() error {
	if r.RequestedAt.IsZero() {
		return ErrInvalidIngestionRequest
	}
	if strings.TrimSpace(r.Reason) == "" {
		return ErrInvalidIngestionRequest
	}
	return nil
}

// IngestionResult is a small summary suitable for printing in runnable examples.
type IngestionResult struct {
	StartedAt  time.Time
	FinishedAt time.Time
	Records    int
	Outcome    string
}

// AlertSignal is emitted by watchdog/SLA examples when a deadline is missed.
type AlertSignal struct {
	At     time.Time
	Reason string
}
