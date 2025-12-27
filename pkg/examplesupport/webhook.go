package examplesupport

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidWebhookRequest = errors.New("invalid webhook request")
)

// WebhookRequest is the JSON payload accepted by the webhook-trigger example.
// See specs/001-trigger-source-examples/contracts/webhook-trigger-api.md.
type WebhookRequest struct {
	EventID        string `json:"event_id,omitempty"`
	SourceCategory string `json:"source_category,omitempty"`
	Note           string `json:"note,omitempty"`
}

func (r WebhookRequest) Validate() error {
	if r.EventID != "" && strings.TrimSpace(r.EventID) == "" {
		return ErrInvalidWebhookRequest
	}
	if r.SourceCategory != "" && !isAllowedWebhookSourceCategory(r.SourceCategory) {
		return ErrInvalidWebhookRequest
	}
	return nil
}

func isAllowedWebhookSourceCategory(cat string) bool {
	switch cat {
	case "files", "object_storage", "database", "warehouse", "rest_api", "observability":
		return true
	default:
		return false
	}
}

// DecodeWebhookRequest reads and parses a webhook request body.
// It returns the parsed request plus the raw request body bytes.
func DecodeWebhookRequest(r io.Reader, maxBytes int64) (WebhookRequest, []byte, error) {
	if r == nil {
		return WebhookRequest{}, nil, ErrInvalidWebhookRequest
	}
	if maxBytes <= 0 {
		maxBytes = 64 * 1024
	}

	lr := &io.LimitedReader{R: r, N: maxBytes + 1}
	b, err := io.ReadAll(lr)
	if err != nil {
		return WebhookRequest{}, nil, err
	}
	if int64(len(b)) > maxBytes {
		return WebhookRequest{}, nil, ErrInvalidWebhookRequest
	}

	var req WebhookRequest
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		return WebhookRequest{}, nil, err
	}
	if err := req.Validate(); err != nil {
		return WebhookRequest{}, nil, err
	}
	return req, b, nil
}

// WebhookHandler implements the local webhook-trigger API.
//
// It validates the request and triggers the example run by calling OnEvent.
// The handler is intended for local development only.
type WebhookHandler struct {
	Now func() time.Time

	// MaxBodyBytes limits request body size.
	MaxBodyBytes int64

	// GenerateEventID is used when the client doesn't provide event_id.
	GenerateEventID func(t time.Time) string

	// OnEvent is invoked after a request is accepted.
	OnEvent func(ctx context.Context, ev TriggerEvent, req WebhookRequest) error
}

type webhookSuccessResponse struct {
	Accepted bool   `json:"accepted"`
	EventID  string `json:"event_id"`
	Message  string `json:"message"`
}

type webhookFailureResponse struct {
	Accepted bool   `json:"accepted"`
	Error    string `json:"error"`
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if w == nil || r == nil {
		return
	}

	if r.Method != http.MethodPost {
		writeWebhookFailure(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" || !strings.Contains(strings.ToLower(contentType), "application/json") {
		writeWebhookFailure(w, http.StatusBadRequest, "content-type must be application/json")
		return
	}

	req, raw, err := DecodeWebhookRequest(r.Body, h.MaxBodyBytes)
	if err != nil {
		writeWebhookFailure(w, http.StatusBadRequest, "invalid request")
		return
	}

	now := h.Now
	if now == nil {
		now = time.Now
	}
	firedAt := now()

	eventID := strings.TrimSpace(req.EventID)
	if eventID == "" {
		gen := h.GenerateEventID
		if gen == nil {
			gen = func(t time.Time) string { return "evt-" + t.UTC().Format("20060102T150405.000000000Z07:00") }
		}
		eventID = gen(firedAt)
	}

	sourceRef := "webhook://local"
	if req.SourceCategory != "" {
		sourceRef = req.SourceCategory
	}

	ev := TriggerEvent{
		Kind:      "webhook",
		Occurred:  firedAt,
		Payload:   raw,
		SourceRef: sourceRef,
	}

	if h.OnEvent != nil {
		if err := h.OnEvent(r.Context(), ev, req); err != nil {
			writeWebhookFailure(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(webhookSuccessResponse{Accepted: true, EventID: eventID, Message: "trigger accepted"})
}

func writeWebhookFailure(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(webhookFailureResponse{Accepted: false, Error: msg})
}
