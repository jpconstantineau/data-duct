# Research: Trigger & Source Examples

**Date**: 2025-12-26

This feature adds runnable examples (no new third-party dependencies) demonstrating trigger patterns and datasource ingestion patterns for the existing pipeline library.

## Decisions

### Decision: Example organization = two sets
- **Chosen**: Two runnable sets
  - Trigger-focused examples (one per trigger style) using simulated/local data
  - Datasource-focused examples (one per datasource category) using a manual/run-now trigger
- **Rationale**: Covers all requirements without a large trigger×datasource matrix; keeps binaries small and runnable without external infrastructure.
- **Alternatives considered**:
  - Matrix pairing (rejected: too many examples)
  - Curated overlap set (rejected: ambiguous coverage)

### Decision: Backpressure when trigger fires during active run
- **Chosen**: Coalesce (queue size 1)
- **Rationale**: Prevents unbounded backlog; avoids concurrent run complexity; still guarantees “run again” after bursts.
- **Alternatives considered**:
  - Skip (may miss data freshness)
  - Queue all (unbounded unless capped)
  - Concurrent runs (harder to reason about and demonstrate)

### Decision: Daily-at-HH:MM scheduling without cron
- **Chosen**: Compute the next occurrence of HH:MM using `time` only.
- **Rationale**: No cron parser dependency required; easy to explain.
- **Notes**:
  - Interpret HH:MM in local time.
  - If the current time is past today’s HH:MM, schedule for tomorrow.
  - Handle clock changes by recomputing next-run time after each run.

### Decision: File “data availability” trigger
- **Chosen**: Polling with `os.Stat` / directory scan on a configurable interval.
- **Rationale**: `fsnotify` would be third-party; polling is dependency-free and portable.
- **Notes**:
  - Default behavior should avoid re-processing the same file: track a small in-memory set of processed file names for the run.

### Decision: Watchdog/SLA trigger semantics
- **Chosen**: Emit an alert event into the pipeline; print an explicit ALERT; exit 0.
- **Rationale**: Keeps behavior testable and visible while following the clarified requirement that watchdog does not fail the process.

### Decision: REST API datasource used for examples
- **Chosen**: `https://uselessfacts.jsph.pl/api/v2/facts/random?language=en`
- **Rationale**: Public endpoint with stable JSON; works with stdlib (`net/http`, `encoding/json`).
- **Notes**:
  - Default example runs “once” to remain deterministic.
  - Use timeouts to avoid hanging.

### Decision: Datasource categories that usually require SDKs/drivers
- **Chosen**: Demonstrate ingestion via dependency-free representations:
  - Databases: ingest exported extracts (CSV/JSON) from disk
  - Warehouses: ingest exported extracts (CSV/JSON) from disk
  - Object storage: local folder stand-in by default
  - Observability: ingest log file lines + metrics over HTTP JSON + trace export file (JSON)
- **Rationale**: Meets “no third-party libraries” while still demonstrating realistic ingestion pathways.

## Open Questions

None. All prior clarifications are resolved in the feature spec.
