# Tasks: Trigger & Source Examples

**Input**: Design documents from `/specs/001-trigger-source-examples/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/, quickstart.md

**Tests**: Tests are REQUIRED. Tasks include test-first work for each story (tests written first, initially failing) and the feature is not complete until all tests pass.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add the minimal shared structure used by all examples.

- [ ] T001 Create example support package skeleton in pkg/examplesupport/doc.go and pkg/examplesupport/types.go
- [ ] T002 Create base fixture folders with placeholders in testdata/trigger-source-examples/exports/.gitkeep and testdata/trigger-source-examples/observability/.gitkeep
- [ ] T003 [P] Add/extend docs entry listing new example commands in docs/pipeline/README.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared helper utilities used across multiple examples.

**⚠️ CRITICAL**: No user story work should start until these helpers are in place (to keep examples consistent and testable).

### Tests (write first)

- [ ] T004 [P] Add unit tests for example event/data types in pkg/examplesupport/types_test.go
- [ ] T005 [P] Add unit tests for coalescing backpressure controller in pkg/examplesupport/coalesce_test.go
- [ ] T006 [P] Add unit tests for shared “process step” helper in pkg/examplesupport/process_test.go

### Implementation

- [ ] T007 Implement shared types (TriggerEvent, IngestionRequest, IngestionResult, AlertSignal) in pkg/examplesupport/types.go
- [ ] T008 Implement coalescing runner (max one pending run) in pkg/examplesupport/coalesce.go
- [ ] T009 Implement a reusable process-step adapter used by trigger examples in pkg/examplesupport/process.go
- [ ] T010 Add a small helper for consistent CLI output formatting in pkg/examplesupport/output.go

**Checkpoint**: Foundation ready — trigger examples can now be implemented consistently.

---

## Phase 3: User Story 1 - Run on a schedule (Priority: P1) 🎯 MVP

**Goal**: Provide runnable examples for interval schedule and daily-at-HH:MM (no cron) that initiate ingestion via a process step.

**Independent Test**: `go test ./...` passes; `go run ./cmd/trigger-interval -interval 200ms -duration 2s` triggers multiple runs; `go run ./cmd/trigger-daily -time <soon> -once` triggers one run.

### Tests for User Story 1 (REQUIRED) ⚠️

- [ ] T011 [P] [US1] Add tests for daily HH:MM next-run calculation in pkg/examplesupport/schedule_test.go
- [ ] T012 [P] [US1] Add tests for interval trigger run loop behavior in pkg/examplesupport/interval_test.go

### Implementation for User Story 1

- [ ] T013 [US1] Implement daily HH:MM next-run calculation in pkg/examplesupport/schedule.go
- [ ] T014 [US1] Implement interval trigger loop helper in pkg/examplesupport/interval.go
- [ ] T015 [US1] Add runnable interval trigger example in cmd/trigger-interval/main.go
- [ ] T016 [US1] Add runnable daily HH:MM trigger example in cmd/trigger-daily/main.go

**Checkpoint**: US1 examples run and demonstrate schedule triggering.

---

## Phase 4: User Story 2 - Run on events (Priority: P2)

**Goal**: Provide runnable examples for webhook/API trigger and manual “run now” trigger.

**Independent Test**: `go test ./...` passes; webhook example accepts a POST and initiates one ingestion run; manual example initiates one ingestion run.

### Tests for User Story 2 (REQUIRED) ⚠️

- [ ] T017 [P] [US2] Add tests for webhook request validation/parsing in pkg/examplesupport/webhook_test.go
- [ ] T018 [P] [US2] Add tests for manual trigger request construction in pkg/examplesupport/manual_test.go

### Implementation for User Story 2

- [ ] T019 [US2] Implement webhook trigger HTTP handler and request parsing in pkg/examplesupport/webhook.go
- [ ] T020 [US2] Implement manual run-now helper in pkg/examplesupport/manual.go
- [ ] T021 [US2] Add runnable webhook trigger example in cmd/trigger-webhook/main.go
- [ ] T022 [US2] Add runnable manual run-now example in cmd/trigger-manual-run-now/main.go

**Checkpoint**: US2 examples run and demonstrate event-driven triggering.

---

## Phase 5: User Story 3 - Detect missing data (Priority: P3)

**Goal**: Provide runnable examples for file-availability polling and watchdog/SLA deadline alerting.

**Independent Test**: `go test ./...` passes; file polling triggers ingestion on new file and does not re-trigger for the same file by default; watchdog emits an ALERT event and exits 0.

### Tests for User Story 3 (REQUIRED) ⚠️

- [ ] T023 [P] [US3] Add tests for file polling “do not reprocess” behavior in pkg/examplesupport/filepoll_test.go
- [ ] T024 [P] [US3] Add tests for watchdog deadline detection and alert emission in pkg/examplesupport/watchdog_test.go

### Implementation for User Story 3

- [ ] T025 [US3] Implement file polling utilities (scan, filter, seen-set) in pkg/examplesupport/filepoll.go
- [ ] T026 [US3] Implement watchdog deadline monitor that emits AlertSignal into pipeline in pkg/examplesupport/watchdog.go
- [ ] T027 [US3] Add runnable file polling trigger example in cmd/trigger-file-poll/main.go
- [ ] T028 [US3] Add runnable watchdog/SLA example in cmd/trigger-watchdog-sla/main.go

**Checkpoint**: US3 examples run and demonstrate polling + alerting.

---

## Phase 6: Datasource-Focused Examples (FR-004)

**Goal**: Provide 5 runnable datasource ingestion examples using a manual/run-now trigger.

**Independent Test**: `go test ./...` passes; each `go run ./cmd/source-*` produces a clear ingestion summary.

### Testdata fixtures (write first where possible)

- [ ] T029 [P] Create DB export fixture in testdata/trigger-source-examples/exports/db-export.csv
- [ ] T030 [P] Create warehouse export fixture in testdata/trigger-source-examples/exports/warehouse-export.csv
- [ ] T031 [P] Create observability log fixture in testdata/trigger-source-examples/observability/app.log
- [ ] T032 [P] Create observability trace export fixture in testdata/trigger-source-examples/observability/traces.json

### Ingestion helper tests (REQUIRED)

- [ ] T033 [P] Add tests for CSV/JSON parsing helpers in pkg/examplesupport/ingest_test.go
- [ ] T034 [P] Add tests for REST uselessfacts response parsing in pkg/examplesupport/rest_uselessfacts_test.go

### Ingestion helper implementation

- [ ] T035 Implement CSV/JSON ingestion helpers in pkg/examplesupport/ingest.go
- [ ] T036 Implement uselessfacts REST fetch + parse helper (stdlib only) in pkg/examplesupport/rest_uselessfacts.go

### Runnable datasource examples

- [ ] T037 Add files/object storage stand-in ingestion example in cmd/source-files/main.go
- [ ] T038 Add database export ingestion example in cmd/source-database-export/main.go
- [ ] T039 Add warehouse export ingestion example in cmd/source-warehouse-export/main.go
- [ ] T040 Add REST API ingestion example using uselessfacts endpoint in cmd/source-rest-uselessfacts/main.go
- [ ] T041 Add observability ingestion example (logs + traces; optional metrics URL) in cmd/source-observability/main.go

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Ensure everything is discoverable, deterministic, and passes quality gates.

- [ ] T042 Update quickstart references (if needed) in specs/001-trigger-source-examples/quickstart.md
- [ ] T048 [P] [US2] Add tests for webhook response contract (status + JSON body) in pkg/examplesupport/webhook_contract_test.go
- [ ] T049 Ensure docs/contracts align with fixtures and commands (update specs/001-trigger-source-examples/quickstart.md and specs/001-trigger-source-examples/contracts/examples-cli.md)
- [ ] T043 [P] Add a short section listing all new commands in docs/pipeline/README.md
- [ ] T044 Run gofmt on pkg/examplesupport/*.go and cmd/*/main.go
- [ ] T045 Run unit tests via `go test` for pkg/examplesupport/* and fix any failures
- [ ] T046 (Optional) Run scripts/ci/lint.ps1 and address any new lint issues
- [ ] T047 (Optional) Run scripts/ci/vuln.ps1 and address any actionable findings

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup; blocks trigger/source work
- **US1/US2/US3**: Depend on Foundational; can be done in priority order or in parallel
- **Datasource Examples (Phase 6)**: Depends on Foundational; can be parallelized per example
- **Polish (Phase 7)**: Depends on all desired examples being complete

### User Story Dependencies

- **US1 (P1)**: No dependencies on US2/US3
- **US2 (P2)**: No dependencies on US1/US3
- **US3 (P3)**: No dependencies on US1/US2

---

## Parallel Execution Examples

### US1 parallel work

- Run in parallel: T011, T012

### US2 parallel work

- Run in parallel: T017, T018

### US3 parallel work

- Run in parallel: T023, T024

### Datasource examples parallel work

- Run in parallel: T037, T038, T039, T040, T041

---

## Implementation Strategy

### MVP First

1. Complete Phase 1–2
2. Complete US1 (interval + daily)
3. Validate with `go test ./...` and run the two US1 examples

### Incremental Delivery

- Add US2, validate
- Add US3, validate
- Add datasource examples, validate
