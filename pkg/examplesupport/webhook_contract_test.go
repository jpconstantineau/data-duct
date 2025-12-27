package examplesupport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type webhookSuccess struct {
	Accepted bool   `json:"accepted"`
	EventID  string `json:"event_id"`
	Message  string `json:"message"`
}

type webhookFailure struct {
	Accepted bool   `json:"accepted"`
	Error    string `json:"error"`
}

func TestWebhookContract_Success202JSON(t *testing.T) {
	h := &WebhookHandler{
		Now:             func() time.Time { return time.Unix(1, 0).UTC() },
		GenerateEventID: func(t time.Time) string { return "evt-test" },
		OnEvent: func(ctx context.Context, ev TriggerEvent, req WebhookRequest) error {
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodPost, "http://example.local/trigger", strings.NewReader(`{"note":"run"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); !strings.Contains(strings.ToLower(ct), "application/json") {
		t.Fatalf("expected application/json content-type, got %q", ct)
	}

	var got webhookSuccess
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !got.Accepted {
		t.Fatalf("expected accepted=true")
	}
	if got.EventID == "" {
		t.Fatalf("expected event_id")
	}
	if got.Message == "" {
		t.Fatalf("expected message")
	}
}

func TestWebhookContract_Validation400JSON(t *testing.T) {
	h := &WebhookHandler{}

	req := httptest.NewRequest(http.MethodPost, "http://example.local/trigger", strings.NewReader(`{"source_category":"nope"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, res.StatusCode)
	}
	var got webhookFailure
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Accepted {
		t.Fatalf("expected accepted=false")
	}
	if got.Error == "" {
		t.Fatalf("expected error message")
	}
}

func TestWebhookContract_Internal500JSON(t *testing.T) {
	h := &WebhookHandler{
		OnEvent: func(ctx context.Context, ev TriggerEvent, req WebhookRequest) error {
			return errTest
		},
	}

	req := httptest.NewRequest(http.MethodPost, "http://example.local/trigger", strings.NewReader(`{"note":"run"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, res.StatusCode)
	}
	var got webhookFailure
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Accepted {
		t.Fatalf("expected accepted=false")
	}
	if got.Error == "" {
		t.Fatalf("expected error message")
	}
}

var errTest = &testError{s: "boom"}

type testError struct{ s string }

func (e *testError) Error() string { return e.s }
