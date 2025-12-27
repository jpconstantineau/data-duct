# data-duct Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-12-24

## Active Technologies
- Go 1.25 + Go standard library only for examples; use existing `pkg/pipeline` library in this repo. (001-trigger-source-examples)
- Local files only (examples are runnable without external services). (001-trigger-source-examples)

- Go 1.25 + None at runtime for core packages (stdlib-only). Optional dev tooling (lint/security) via CI. (001-graceful-context-pipeline)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.25

## Code Style

Go 1.25: Follow standard conventions

## Recent Changes
- 001-trigger-source-examples: Added Go 1.25 + Go standard library only for examples; use existing `pkg/pipeline` library in this repo.

- 001-graceful-context-pipeline: Added Go 1.25 + None at runtime for core packages (stdlib-only). Optional dev tooling (lint/security) via CI.

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->

## Constitution

All development guidelines and generated artifacts MUST comply with
`.specify/memory/constitution.md`.
