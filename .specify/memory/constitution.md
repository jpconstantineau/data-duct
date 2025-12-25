<!--
Sync Impact Report

- Version change: (template) -> 0.1.0
- Modified principles:
	- Principle 1 -> Library-First & Reusable
	- Principle 2 -> Test-First (TDD) (NON-NEGOTIABLE)
	- Principle 3 -> Core Has No Runtime Dependencies
	- Principle 4 -> Spec-Kit-Driven Development
	- Principle 5 -> Quality Gates & Clean Architecture
- Added sections: Additional Constraints, Development Workflow
- Removed sections: none
- Templates requiring updates:
	- .specify/templates/plan-template.md ✅ updated
	- .specify/templates/spec-template.md ✅ updated
	- .specify/templates/tasks-template.md ✅ updated
	- .specify/templates/checklist-template.md ✅ updated
	- .specify/templates/agent-file-template.md ✅ updated
- Follow-up TODOs:
	- RATIFICATION_DATE: TODO (original adoption date unknown)
-->

# data-duct Constitution

## Core Principles

### Library-First & Reusable
Every feature MUST begin as a standalone, reusable library (a package/module that can be used
independently). Features MUST avoid entangling business logic with integrations; integration
adapters are separate packages.

Each feature/library MUST include:

- A minimal public API with clear responsibilities
- Complete tests
- Developer documentation
- At least one runnable example demonstrating intended usage

### Test-First (TDD) (NON-NEGOTIABLE)
Development MUST be test-first:

- Tests MUST be written before implementation.
- Tests MUST fail before implementation is written.
- A feature is NOT complete until all tests are complete and passing.

### Core Has No Runtime Dependencies
The `core` library MUST have no runtime dependencies beyond the Go standard library.

Practically, this means:

- `core` packages MUST NOT import third-party modules.
- Any integrations (cloud SDKs, DB clients, queue clients, etc.) MUST live outside `core`.
- `core` MUST define interfaces/contracts that adapters implement.

### Spec-Kit-Driven Development
All new features MUST be developed using the spec-kit methodology.

Feature work MUST be traceable:

`specs/[###-feature-name]/spec.md` → `plan.md` → `tasks.md` → implementation + tests + examples

### Quality Gates & Clean Architecture
The library MUST follow clean code and clean architecture practices.

All new/changed code MUST pass:

- Code formatting (Go: `gofmt` and typically `goimports`)
- Static analysis (e.g., `golangci-lint`)
- Security checks (e.g., `govulncheck`)
- Unit tests (`go test ./...`)
- Code coverage checks (thresholds defined by project governance)

The repository MUST be organized according to golang-standards/project-layout.

## Additional Constraints

- Developer-facing documentation MUST be kept current and actionable.
- Each feature MUST include examples showing how to use it.
- Keep the core small and dependency-free; push integrations to adapters.

## Development Workflow

- All changes MUST be driven by spec-kit artifacts (spec → plan → tasks).
- PRs MUST demonstrate test-first development (tests added first, then implementation).
- CI MUST enforce formatting, linting/static analysis, security checks, unit tests, and
	coverage gates.

## Governance

This constitution supersedes other conventions when conflicts exist.

Amendments

- Changes MUST be proposed via PR with rationale and migration plan (if needed).
- `CONSTITUTION_VERSION` follows semantic versioning:
	- MAJOR: backward-incompatible governance change or removal/redefinition of principles
	- MINOR: new principle/section or material expansion
	- PATCH: clarifications/wording/typos without semantic change

Compliance

- Feature plans MUST include a "Constitution Check" section based on this document.
- Reviews MUST verify compliance, and document any approved exceptions.

**Version**: 0.1.0 | **Ratified**: TODO(RATIFICATION_DATE): original adoption date unknown | **Last Amended**: 2025-12-24
