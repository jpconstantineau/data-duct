# Contract: Example CLIs

**Date**: 2025-12-26

This contract lists expected command-line options for runnable examples.

## Global conventions

- All examples are runnable via `go run ./cmd/<name>`.
- Flags use Go `flag` package.
- Output must clearly indicate: trigger fired, ingestion started, ingestion finished, outcome.

## Trigger-focused examples

### Interval schedule
- Flags:
  - `-interval` (default e.g. `1s`)
  - `-duration` (optional; if set, example exits after duration)

### Daily at HH:MM
- Flags:
  - `-time` (string `HH:MM`)
  - `-once` (bool; default true for deterministic behavior)

### Webhook/API trigger
- Flags:
  - `-listen` (default e.g. `127.0.0.1:8080`)
  - `-path` (default `/trigger`)

### File-availability polling
- Flags:
  - `-watch` (directory)
  - `-pattern` (optional; matched against the base filename using Go `filepath.Match`)
  - `-poll` (duration)

### Manual trigger (run now)
- Flags:
  - `-run-now` (bool)

### Watchdog/SLA
- Flags:
  - `-deadline` (duration)
  - `-poll` (duration)

## Datasource-focused examples (manual/run-now)

### Files / object storage (local folder stand-in)
- Flags:
  - `-root` (directory)

### Database (dependency-free)
- Flags:
  - `-input` (path to exported data file, e.g. CSV/JSON)
  - Example fixture path: `./testdata/trigger-source-examples/exports/db-export.csv`

### Data warehouse (dependency-free)
- Flags:
  - `-input` (path to exported data file, e.g. CSV/JSON)
  - Example fixture path: `./testdata/trigger-source-examples/exports/warehouse-export.csv`

### REST API (uselessfacts)
- Flags:
  - `-url` (default `https://uselessfacts.jsph.pl/api/v2/facts/random?language=en`)

### Observability (metrics/traces/logs)
- Flags:
  - `-logs` (path to log file)
  - `-traces` (path to trace export file)
  - `-metrics-url` (URL to pull metrics JSON, optional)
  - Example fixture paths:
    - logs: `./testdata/trigger-source-examples/observability/app.log`
    - traces: `./testdata/trigger-source-examples/observability/traces.json`
