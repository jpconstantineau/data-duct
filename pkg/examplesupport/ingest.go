package examplesupport

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
)

var ErrInvalidInput = errors.New("invalid input")

// CountCSVRecords counts CSV records from r.
// If hasHeader is true, the first record is excluded from the count.
func CountCSVRecords(r io.Reader, hasHeader bool) (int, error) {
	if r == nil {
		return 0, ErrInvalidInput
	}
	cr := csv.NewReader(r)
	recs, err := cr.ReadAll()
	if err != nil {
		return 0, err
	}
	if len(recs) == 0 {
		return 0, nil
	}
	count := len(recs)
	if hasHeader && count > 0 {
		count--
	}
	if count < 0 {
		count = 0
	}
	return count, nil
}

// CountLogLines counts newline-delimited lines in r.
func CountLogLines(r io.Reader) (int, error) {
	if r == nil {
		return 0, ErrInvalidInput
	}
	s := bufio.NewScanner(r)
	// Allow longer lines than the default 64K, to avoid surprising failures.
	s.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	lines := 0
	for s.Scan() {
		lines++
	}
	if err := s.Err(); err != nil {
		return 0, err
	}
	return lines, nil
}

// CountJSONArrayItems counts the number of items in a top-level JSON array.
func CountJSONArrayItems(r io.Reader) (int, error) {
	if r == nil {
		return 0, ErrInvalidInput
	}
	dec := json.NewDecoder(r)

	tok, err := dec.Token()
	if err != nil {
		return 0, err
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '[' {
		return 0, ErrInvalidInput
	}

	count := 0
	for dec.More() {
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return 0, err
		}
		count++
	}

	endTok, err := dec.Token()
	if err != nil {
		return 0, err
	}
	endDelim, ok := endTok.(json.Delim)
	if !ok || endDelim != ']' {
		return 0, ErrInvalidInput
	}

	return count, nil
}
