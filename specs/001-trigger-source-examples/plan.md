# Implementation Plan: Trigger & Source Examples

**Branch**: `001-trigger-source-examples` | **Date**: 2025-12-26 | **Spec**: `specs/001-trigger-source-examples/spec.md`
**Input**: Feature specification from `/specs/001-trigger-source-examples/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Add runnable Go examples demonstrating:

- Trigger styles: interval schedule, daily at HH:MM (no cron), webhook/API event, file-availability polling, manual run-now, watchdog/SLA deadline
- Datasource ingestion styles: files/object storage (local folder stand-in), database exports, data warehouse exports, REST API (using the uselessfacts endpoint), observability (metrics/traces/logs via dependency-free inputs)

Examples must be standard-library-only (no third-party dependencies), Windows-friendly, and use the existing pipeline library (`pkg/pipeline`) by sending trigger events to a “process” step that performs the ingestion.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.25
**Primary Dependencies**: Go standard library only for examples; use existing `pkg/pipeline` library in this repo.
**Storage**: Local files only (examples are runnable without external services).
**Testing**: `go test ./...` (unit tests for helper packages and lightweight integration tests for examples where feasible)
**Target Platform**: Windows (primary), but keep examples cross-platform.
**Project Type**: Library repo with runnable `cmd/` examples (golang-standards/project-layout)
**Performance Goals**: Not performance-focused; correctness and clarity of example behavior.
**Constraints**:
- No third-party dependencies
- Deterministic default runs (examples should finish quickly by default)
- Backpressure default: coalesce (max one pending run)
- Watchdog “deadline missed” emits alert event and exits 0
**Scale/Scope**: 6 trigger-focused runnable examples + 5 datasource-focused runnable examples.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Plans MUST include a short "Constitution Compliance" statement derived from
`.specify/memory/constitution.md`, covering:

- Library-first packaging and reusable boundaries
- Test-first approach (tests written first, initially failing)
- `core` dependency constraints (stdlib-only imports)
- Quality gates (formatting, static analysis, security checks, tests, coverage)
- Example(s) included for the feature

Constitution Compliance (this feature)

- Library-first: introduce a small reusable helper library `pkg/examplesupport` (stdlib-only) that encapsulates trigger orchestration primitives and ingestion helpers; runnable `cmd/` examples consume this library alongside `pkg/pipeline`.
- Test-first: any new helper packages introduced to support examples will be driven by tests first.
- Core constraints: no new runtime dependencies are introduced; existing core packages remain stdlib-only.
- Quality gates: gofmt + existing CI scripts (test/lint/vuln) continue to pass.
- Examples: multiple runnable example commands are added under `cmd/`.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
cmd/
├── trigger-interval/            Interval schedule trigger → process step ingests simulated/local data
├── trigger-daily/               Daily HH:MM trigger (no cron) → process step ingests simulated/local data
├── trigger-webhook/             HTTP webhook trigger → process step ingests simulated/local data
├── trigger-file-poll/           File availability polling trigger → process step ingests simulated/local data
├── trigger-manual-run-now/      Manual run-now trigger → process step ingests simulated/local data
└── trigger-watchdog-sla/        Watchdog trigger emits alert event when deadline missed

cmd/
├── source-files/                Manual trigger → process step ingests files from a folder (also serves as object storage stand-in)
├── source-database-export/      Manual trigger → process step ingests exported DB extract (CSV/JSON)
├── source-warehouse-export/     Manual trigger → process step ingests exported warehouse extract (CSV/JSON)
├── source-rest-uselessfacts/    Manual trigger → process step fetches REST data from uselessfacts API
└── source-observability/        Manual trigger → process step ingests logs/metrics/traces from dependency-free inputs

pkg/
└── pipeline/                    Existing public library API used by all examples

pkg/
└── examplesupport/              Reusable helper library for trigger orchestration + ingestion helpers (stdlib-only)

testdata/
└── (optional)                  Small local fixtures for exports/logs/traces
```

**Structure Decision**: Add runnable example binaries under `cmd/` and a small reusable helper library under `pkg/examplesupport` (stdlib-only) plus minimal `testdata/` fixtures.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |

## Phase 0: Outline & Research (complete)

- Outputs:
  - `specs/001-trigger-source-examples/research.md`
- Notes:
  - All clarifications are resolved in the feature spec.
  - REST API chosen for test ingestion: `https://uselessfacts.jsph.pl/api/v2/facts/random?language=en`.

## Phase 1: Design & Contracts (complete)

- Outputs:
  - `specs/001-trigger-source-examples/data-model.md`
  - `specs/001-trigger-source-examples/contracts/examples-cli.md`
  - `specs/001-trigger-source-examples/contracts/webhook-trigger-api.md`
  - `specs/001-trigger-source-examples/contracts/rest-uselessfacts-api.md`
  - `specs/001-trigger-source-examples/quickstart.md`

## Phase 1.5: Agent Context Update

- Run:
  - `.specify/scripts/powershell/update-agent-context.ps1 -AgentType copilot`

## Phase 2: Implementation Planning (to be executed next)

1. Add reusable helper library package(s) in `pkg/examplesupport` for:
   - backpressure coalescing controller
   - daily schedule next-run calculation
   - file polling state (seen/processed files)
2. Add trigger-focused example commands under `cmd/` (6 total).
3. Add datasource-focused example commands under `cmd/` (5 total), including REST ingestion using `uselessfacts`.
4. Add minimal `testdata/` fixtures for DB/warehouse exports and observability logs/traces.
5. Add tests:
   - unit tests for helpers (TDD)
   - small “smoke” tests where practical (avoid flaky timing dependencies)
6. Validate quality gates:
   - `gofmt ./...`
   - `go test ./...`
   - optional: `scripts/ci/lint.ps1` and `scripts/ci/vuln.ps1`
