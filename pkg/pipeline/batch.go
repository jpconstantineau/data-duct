package pipeline

import "time"

// BatchPolicy controls how a batch stage groups items.
type BatchPolicy struct {
	Size    int
	MaxWait time.Duration
}
