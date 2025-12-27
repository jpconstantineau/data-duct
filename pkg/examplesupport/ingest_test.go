package examplesupport

import (
	"strings"
	"testing"
)

func TestCountCSVRecords_ExcludesHeader(t *testing.T) {
	csv := "id,name\n1,alpha\n2,beta\n"
	got, err := CountCSVRecords(strings.NewReader(csv), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCountLogLines_CountsLines(t *testing.T) {
	log := "a\n\nccc\n"
	got, err := CountLogLines(strings.NewReader(log))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestCountJSONArrayItems_CountsObjects(t *testing.T) {
	json := `[{"a":1},{"b":2},{"c":3}]`
	got, err := CountJSONArrayItems(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}
