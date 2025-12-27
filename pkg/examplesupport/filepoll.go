package examplesupport

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

// FilePollConfig configures the file-availability polling helper.
//
// Pattern, when set, is matched against the base filename using filepath.Match.
// Seen tracking is in-memory; it prevents re-ingesting the same file by default.
type FilePollConfig struct {
	WatchDir string
	Pattern  string
	Now      func() time.Time

	ReadDir func(name string) ([]os.DirEntry, error)
}

// FilePoller scans a directory and emits TriggerEvents for new matching files.
//
// It is intentionally small, stdlib-only, and designed for deterministic unit tests.
type FilePoller struct {
	watchDir string
	pattern  string
	now      func() time.Time
	readDir  func(string) ([]os.DirEntry, error)

	seen map[string]struct{}
}

func NewFilePoller(cfg FilePollConfig) *FilePoller {
	now := cfg.Now
	if now == nil {
		now = time.Now
	}
	readDir := cfg.ReadDir
	if readDir == nil {
		readDir = os.ReadDir
	}

	return &FilePoller{
		watchDir: cfg.WatchDir,
		pattern:  cfg.Pattern,
		now:      now,
		readDir:  readDir,
		seen:     make(map[string]struct{}),
	}
}

// Scan returns the next new matching file, if any.
//
// The returned TriggerEvent has Kind "file_poll" and SourceRef set to the full file path.
func (p *FilePoller) Scan() (TriggerEvent, bool, error) {
	if p == nil {
		return TriggerEvent{}, false, nil
	}
	if p.watchDir == "" {
		return TriggerEvent{}, false, nil
	}

	entries, err := p.readDir(p.watchDir)
	if err != nil {
		return TriggerEvent{}, false, err
	}

	candidates := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if _, ok := p.seen[name]; ok {
			continue
		}
		if p.pattern != "" {
			matched, matchErr := filepath.Match(p.pattern, filepath.Base(name))
			if matchErr != nil {
				return TriggerEvent{}, false, matchErr
			}
			if !matched {
				continue
			}
		}
		candidates = append(candidates, name)
	}

	if len(candidates) == 0 {
		return TriggerEvent{}, false, nil
	}

	sort.Strings(candidates)
	picked := candidates[0]
	p.seen[picked] = struct{}{}

	full := filepath.Join(p.watchDir, picked)
	return TriggerEvent{
		Kind:      "file_poll",
		Occurred:  p.now(),
		SourceRef: full,
	}, true, nil
}
