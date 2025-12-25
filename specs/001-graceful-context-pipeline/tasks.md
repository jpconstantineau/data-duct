---

description: "Task list for Graceful Context Pipeline"
---

# Tasks: Graceful Context Pipeline

**Input**: Design documents from `/specs/001-graceful-context-pipeline/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/, quickstart.md

**Tests**: Tests are REQUIRED. Tasks MUST include test-first work for each story (tests written first, initially failing) and the feature is not complete until all tests pass.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `- [ ] T### [P?] [US#?] Description with file path`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[US#]**: Which user story this task belongs to (US1, US2, US3)
- Every task includes an explicit file path

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize Go module, repo structure, and quality gates

- [x] T001 Create project-layout directories in cmd/, pkg/, internal/, docs/, test/ (repo root)
- [x] T002 Initialize Go module for repo in go.mod (repo root)
- [x] T003 [P] Add baseline README for pipeline feature in docs/pipeline/README.md
- [x] T004 [P] Add CI task runner scripts in scripts/ci/{fmt.ps1,lint.ps1,test.ps1,vuln.ps1}
- [x] T005 [P] Add GitHub Actions workflow for gofmt/go test/coverage/lint/vulncheck in .github/workflows/ci.yml
- [x] T006 [P] Add golangci-lint configuration in .golangci.yml

**Checkpoint**: `go test ./...` runs (even if no packages yet) and CI pipeline skeleton exists

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core public API skeleton + internal wiring boundaries (no business logic)

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T007 Define public handler and feed types in pkg/pipeline/types.go
- [x] T008 Define public options and stage options in pkg/pipeline/options.go
- [x] T009 Define public result types (Succeeded/Cancelled/Failed) in pkg/pipeline/result.go
- [x] T010 Define internal stage/wiring primitives (worker orchestration, channels, shutdown) in internal/pipelineinternal/wiring.go
- [x] T011 Define internal error/cancel policy implementation in internal/pipelineinternal/policy.go
- [x] T012 [P] Add package-level doc comments for pkg/pipeline in pkg/pipeline/doc.go

**Checkpoint**: Public API compiles, `pkg/pipeline` imports only stdlib, and internal package is not imported by external users

---

## Phase 3: User Story 1 - Minimal pipeline (Priority: P1) 🎯 MVP

**Goal**: Build and run a minimal source → processor → sink pipeline with compile-time type safety

**Independent Test**: Deterministic in-memory source produces known items; processor transforms; sink receives expected outputs

### Tests for User Story 1 (REQUIRED) ⚠️

> NOTE: Write these tests FIRST, ensure they FAIL before implementation

- [x] T013 [P] [US1] Add compile-time API contract test (type chaining) in pkg/pipeline/contract_test.go
- [x] T014 [P] [US1] Add end-to-end happy-path test (source→then→to) in pkg/pipeline/pipeline_happy_test.go
- [x] T015 [P] [US1] Add processor error propagation test in pkg/pipeline/pipeline_error_test.go

### Implementation for User Story 1

- [x] T016 [US1] Implement Pipeline builder `New` and `Then` chaining in pkg/pipeline/pipeline.go
- [x] T017 [US1] Implement `To` and `Run` entrypoint in pkg/pipeline/runnable.go
- [x] T018 [US1] Implement internal worker loop for single-item stages in internal/pipelineinternal/worker_single.go
- [x] T019 [US1] Implement source pump (read from `<-chan T` until closed/cancelled) in internal/pipelineinternal/source.go
- [x] T020 [US1] Implement sink consumer loop in internal/pipelineinternal/sink.go

### Example (REQUIRED)

- [x] T021 [US1] Add runnable example binary in cmd/graceful-context-pipeline-example/main.go

**Checkpoint**: US1 tests pass and the example runs successfully

---

## Phase 4: User Story 2 - Graceful shutdown & cancellation (Priority: P2)

**Goal**: Respect root context cancellation and ensure deterministic shutdown with no goroutine leaks

**Independent Test**: Run pipeline with a blocking processor/sink and cancel context; `Run` returns promptly and internal goroutines complete

### Tests for User Story 2 (REQUIRED) ⚠️

- [x] T022 [P] [US2] Add cancellation test (cancel while running) in pkg/pipeline/pipeline_cancel_test.go
- [x] T023 [P] [US2] Add shutdown completion test (internal waitgroups complete) in pkg/pipeline/pipeline_shutdown_test.go
- [x] T024 [P] [US2] Add panic-in-handler policy test (panic becomes error) in pkg/pipeline/pipeline_panic_test.go

### Implementation for User Story 2

- [x] T025 [US2] Implement cancellation propagation and upstream stop behavior in internal/pipelineinternal/cancel.go
- [x] T026 [US2] Implement best-effort in-flight completion bounded by ctx in internal/pipelineinternal/drain.go
- [x] T027 [US2] Implement panic recovery wrapper around user handlers in internal/pipelineinternal/safehandler.go
- [x] T028 [US2] Ensure deterministic channel close ordering (no sends on closed channels) in internal/pipelineinternal/wiring.go

**Checkpoint**: US2 tests pass and cancellation behavior matches spec clarifications

---

## Phase 5: User Story 3 - Optional batching (Priority: P3)

**Goal**: Allow a batch processing stage using `BatchHandler` with deterministic batch boundaries

**Independent Test**: Configure a batch stage with fixed size; verify handler receives expected groupings and outputs forward correctly

### Tests for User Story 3 (REQUIRED) ⚠️

- [x] T029 [P] [US3] Add batch grouping behavior test in pkg/pipeline/pipeline_batch_test.go
- [x] T030 [P] [US3] Add batch flush-on-close/cancel test in pkg/pipeline/pipeline_batch_flush_test.go

### Implementation for User Story 3

- [x] T031 [US3] Define batch policy (fixed size, optional max wait) in pkg/pipeline/batch.go
- [x] T032 [US3] Add `ThenBatch` public API in pkg/pipeline/pipeline_batch.go
- [x] T033 [US3] Implement internal batch stage worker in internal/pipelineinternal/worker_batch.go

**Checkpoint**: US3 tests pass and batching can be used alongside single-item stages

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, quality gates, and hardening

- [x] T034 [P] Expand developer documentation (usage, cancellation semantics, error policy) in docs/pipeline/README.md
- [x] T035 [P] Add package examples in pkg/pipeline/example_test.go (Go doc examples)
- [x] T036 Add structured logging hooks (stdlib `log/slog`) behind options in pkg/pipeline/logging.go
- [x] T037 [P] Add security gate documentation and commands to README in docs/pipeline/README.md
- [x] T038 Run quickstart validation updates in specs/001-graceful-context-pipeline/quickstart.md


---

## Dependencies & Execution Order

### User Story dependency graph

- Setup (Phase 1) → Foundational (Phase 2) → US1 (Phase 3)
- US2 depends on US1
- US3 depends on US1
- Polish depends on all selected user stories

### Parallel opportunities

- Phase 1: T003–T006 can be parallelized
- Phase 2: T007–T012 can be partially parallelized (types/options/result/docs)
- US1: tests (T013–T015) parallel; implementation tasks can be split by pkg vs internal
- US2/US3: tests and internal worker implementation can proceed in parallel after US1 wiring exists

---

## Parallel Example: User Story 1

- Run in parallel:
  - T013 [P] [US1] pkg/pipeline/contract_test.go
  - T014 [P] [US1] pkg/pipeline/pipeline_happy_test.go
  - T015 [P] [US1] pkg/pipeline/pipeline_error_test.go

---

## Implementation Strategy

- MVP first: complete Phase 1 → Phase 2 → Phase 3 (US1) and stop to validate example + tests
- Add cancellation/shutdown (US2), then batching (US3)
- Keep `pkg/pipeline` stdlib-only; any optional integrations stay outside core
