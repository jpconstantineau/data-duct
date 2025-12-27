package examplesupport

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFilePoller_Scan_DoesNotReprocessSameFile(t *testing.T) {
	dir := t.TempDir()

	p := NewFilePoller(FilePollConfig{
		WatchDir: dir,
		Now:      func() time.Time { return time.Unix(10, 0).UTC() },
	})

	// Create file after poller is created (doesn't matter; scan reads current dir).
	file1 := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(file1, []byte("a"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	ev, ok, err := p.Scan()
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if ev.Kind != "file_poll" {
		t.Fatalf("expected kind file_poll, got %q", ev.Kind)
	}
	if ev.SourceRef != file1 {
		t.Fatalf("expected sourceRef %q, got %q", file1, ev.SourceRef)
	}

	// Second scan should not return the same file.
	ev2, ok2, err := p.Scan()
	if err != nil {
		t.Fatalf("scan2: %v", err)
	}
	if ok2 {
		t.Fatalf("expected ok=false (no new files), got event %v", ev2)
	}
}

func TestFilePoller_Scan_RespectsPatternOnBaseName(t *testing.T) {
	dir := t.TempDir()

	p := NewFilePoller(FilePollConfig{
		WatchDir: dir,
		Pattern:  "*.csv",
		Now:      func() time.Time { return time.Unix(10, 0).UTC() },
	})

	csv := filepath.Join(dir, "data.csv")
	txt := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(txt, []byte("x"), 0o644); err != nil {
		t.Fatalf("write txt: %v", err)
	}
	if err := os.WriteFile(csv, []byte("c"), 0o644); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	ev, ok, err := p.Scan()
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if ev.SourceRef != csv {
		t.Fatalf("expected %q, got %q", csv, ev.SourceRef)
	}

	_, ok2, err := p.Scan()
	if err != nil {
		t.Fatalf("scan2: %v", err)
	}
	if ok2 {
		t.Fatalf("expected no more matches")
	}
}
