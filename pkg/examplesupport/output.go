package examplesupport

import (
	"fmt"
	"io"
	"time"
)

// PrintOutcome emits consistent, human-readable output for runnable examples.
func PrintOutcome(w io.Writer, example string, outcome string, summary string) {
	if w == nil {
		return
	}
	fmt.Fprintf(w, "[%s] outcome=%s\n", example, outcome)
	if summary != "" {
		fmt.Fprintf(w, "[%s] summary=%s\n", example, summary)
	}
}

// FormatTime is a small helper to keep timestamps consistent across examples.
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return "(zero)"
	}
	return t.Format(time.RFC3339)
}
