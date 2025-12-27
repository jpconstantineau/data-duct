package examplesupport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestParseUselessFactText_OK(t *testing.T) {
	body := []byte(`{"text":"hello","id":"x"}`)
	text, err := ParseUselessFactText(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "hello" {
		t.Fatalf("expected hello, got %q", text)
	}
}

func TestFetchUselessFactText_UsesHTTP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"text":"from-server"}`))
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	text, err := FetchUselessFactText(ctx, srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "from-server" {
		t.Fatalf("expected from-server, got %q", text)
	}
}
