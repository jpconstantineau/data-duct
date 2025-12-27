package examplesupport

import (
	"strings"
	"testing"
)

func TestWebhookRequestValidate_AllowsEmptyOptionalFields(t *testing.T) {
	req := WebhookRequest{}
	if err := req.Validate(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestWebhookRequestValidate_RejectsPresentButEmptyEventID(t *testing.T) {
	req := WebhookRequest{EventID: "   "}
	if err := req.Validate(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestWebhookRequestValidate_RejectsInvalidSourceCategory(t *testing.T) {
	req := WebhookRequest{SourceCategory: "not-a-category"}
	if err := req.Validate(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestDecodeWebhookRequest_RejectsInvalidJSON(t *testing.T) {
	_, _, err := DecodeWebhookRequest(strings.NewReader("{not-json}"), 64*1024)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestDecodeWebhookRequest_ParsesValidJSON(t *testing.T) {
	body := `{"event_id":"e1","source_category":"database","note":"run"}`
	req, raw, err := DecodeWebhookRequest(strings.NewReader(body), 64*1024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.EventID != "e1" {
		t.Fatalf("expected event_id e1, got %q", req.EventID)
	}
	if req.SourceCategory != "database" {
		t.Fatalf("expected source_category database, got %q", req.SourceCategory)
	}
	if req.Note != "run" {
		t.Fatalf("expected note run, got %q", req.Note)
	}
	if string(raw) != body {
		t.Fatalf("expected raw body preserved")
	}
}
