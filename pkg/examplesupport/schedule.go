package examplesupport

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidHHMM = errors.New("invalid HH:MM")

// ParseHHMM parses a string formatted as HH:MM (zero-padded) into hour/minute.
func ParseHHMM(s string) (hour int, minute int, err error) {
	if len(s) != len("00:00") {
		return 0, 0, ErrInvalidHHMM
	}
	if s[2] != ':' {
		return 0, 0, ErrInvalidHHMM
	}
	if s[0] < '0' || s[0] > '9' || s[1] < '0' || s[1] > '9' || s[3] < '0' || s[3] > '9' || s[4] < '0' || s[4] > '9' {
		return 0, 0, ErrInvalidHHMM
	}

	h, err := strconv.Atoi(s[:2])
	if err != nil {
		return 0, 0, ErrInvalidHHMM
	}
	m, err := strconv.Atoi(s[3:])
	if err != nil {
		return 0, 0, ErrInvalidHHMM
	}
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, 0, ErrInvalidHHMM
	}
	return h, m, nil
}

// NextDailyAt returns the next occurrence of hour:minute strictly after now.
//
// If the target time is not after now (equal or before), the next day is returned.
func NextDailyAt(now time.Time, hour int, minute int) time.Time {
	loc := now.Location()
	candidate := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, loc)
	if candidate.After(now) {
		return candidate
	}

	tomorrow := now.AddDate(0, 0, 1)
	return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), hour, minute, 0, 0, loc)
}

// FormatHHMM is a small helper for printing.
func FormatHHMM(hour int, minute int) string {
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

// MustParseHHMM is a convenience for examples; it panics on invalid input.
func MustParseHHMM(s string) (hour int, minute int) {
	h, m, err := ParseHHMM(strings.TrimSpace(s))
	if err != nil {
		panic(err)
	}
	return h, m
}
