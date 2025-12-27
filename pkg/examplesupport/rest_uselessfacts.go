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

var ErrInvalidUselessFactsResponse = errors.New("invalid uselessfacts response")

type uselessFactResponse struct {
	Text string `json:"text"`
}

// ParseUselessFactText parses the uselessfacts response JSON body and returns the fact text.
func ParseUselessFactText(body []byte) (string, error) {
	var resp uselessFactResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	text := strings.TrimSpace(resp.Text)
	if text == "" {
		return "", ErrInvalidUselessFactsResponse
	}
	return text, nil
}

// FetchUselessFactText fetches JSON from url and returns the "text" field.
//
// The request uses a small timeout to avoid hanging in examples.
func FetchUselessFactText(ctx context.Context, url string) (string, error) {
	if strings.TrimSpace(url) == "" {
		return "", ErrInvalidInput
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", ErrInvalidUselessFactsResponse
	}

	lr := &io.LimitedReader{R: res.Body, N: 1024 * 1024}
	b, err := io.ReadAll(lr)
	if err != nil {
		return "", err
	}

	return ParseUselessFactText(b)
}
