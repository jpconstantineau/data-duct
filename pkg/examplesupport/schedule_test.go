package examplesupport

import (
	"testing"
	"time"
)

func TestParseHHMM(t *testing.T) {
	gotH, gotM, err := ParseHHMM("09:30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotH != 9 || gotM != 30 {
		t.Fatalf("expected 9:30, got %d:%d", gotH, gotM)
	}

	_, _, err = ParseHHMM("")
	if err == nil {
		t.Fatalf("expected error for empty")
	}
	_, _, err = ParseHHMM("9:30")
	if err == nil {
		t.Fatalf("expected error for non-zero-padded hour")
	}
	_, _, err = ParseHHMM("09:3")
	if err == nil {
		t.Fatalf("expected error for non-zero-padded minute")
	}
	_, _, err = ParseHHMM("24:00")
	if err == nil {
		t.Fatalf("expected error for hour out of range")
	}
	_, _, err = ParseHHMM("23:60")
	if err == nil {
		t.Fatalf("expected error for minute out of range")
	}
}

func TestNextDailyAt(t *testing.T) {
	loc := time.FixedZone("X", 0)
	now := time.Date(2025, 12, 26, 10, 0, 0, 0, loc)

	next := NextDailyAt(now, 10, 5)
	want := time.Date(2025, 12, 26, 10, 5, 0, 0, loc)
	if !next.Equal(want) {
		t.Fatalf("expected %v, got %v", want, next)
	}

	// If the target time is not after now, schedule for the next day.
	now2 := time.Date(2025, 12, 26, 10, 5, 0, 0, loc)
	next2 := NextDailyAt(now2, 10, 5)
	want2 := time.Date(2025, 12, 27, 10, 5, 0, 0, loc)
	if !next2.Equal(want2) {
		t.Fatalf("expected %v, got %v", want2, next2)
	}
}
