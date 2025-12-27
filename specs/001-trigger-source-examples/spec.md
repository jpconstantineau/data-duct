# Feature Specification: Trigger & Source Examples

**Feature Branch**: `001-trigger-source-examples`  
**Created**: 2025-12-26  
**Status**: Draft  
**Input**: User description: "Create runnable examples for pipeline triggers (interval schedule, daily at HH:MM without cron, webhook/API event, file-availability polling, manual run-now, watchdog/SLA deadline) that send a message to a process step which connects to a datasource and ingests data; include examples for files/object storage, databases, data warehouses, REST API, and observability metrics/traces/logs; do not include examples that require third-party libraries."

## Clarifications

### Session 2025-12-26

- Q: Which default behavior should the trigger examples implement when a trigger fires while an ingestion run is already in progress? → A: Coalesce (allow at most one pending run)
- Q: For the “object storage” datasource example, which dependency-free access pattern should the examples use by default? → A: Local folder as a stand-in
- Q: For the watchdog/SLA trigger example, how should a missed deadline be signaled? → A: Emit an “alert signal” event into the pipeline and exit 0
- Q: How should the runnable examples be organized to satisfy both the trigger list and the datasource list? → A: Two sets (one per trigger using simulated sources, and one per datasource using manual/run-now)

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - Run on a schedule (Priority: P1)

As a developer evaluating the pipeline library, I want runnable examples that trigger ingestion on a fixed interval and on a specific daily time, so I can see how to build time-based automation without relying on external scheduling dependencies.

**Why this priority**: Time-based execution is the most common trigger style and is the fastest way to verify end-to-end behavior.

**Independent Test**: Can be fully tested by running the two schedule examples and observing that ingestion runs at the expected times (interval and daily HH:MM) and reports a successful completion.

**Acceptance Scenarios**:

1. **Given** a configured interval $N$, **When** the example is started, **Then** ingestion is initiated every $N$ interval until stopped.
2. **Given** a configured daily time $HH:MM$, **When** the example is started before that time, **Then** ingestion initiates at the next occurrence of $HH:MM$ without using a cron parser.

---

### User Story 2 - Run on events (Priority: P2)

As a developer building event-driven ingestion, I want runnable examples that trigger ingestion on an incoming webhook/API call and on a manual “run now” request, so I can wire external events and operator actions into the same ingestion pathway.

**Why this priority**: Event-driven ingestion is a key pattern for real-time systems and operational workflows.

**Independent Test**: Can be fully tested by starting the webhook example, sending a request, and observing that a single ingestion run is initiated; and by running the manual trigger example and verifying it performs a single ingestion run.

**Acceptance Scenarios**:

1. **Given** an incoming webhook request with a valid payload, **When** the request is received, **Then** the pipeline starts one ingestion run using details derived from the event.
2. **Given** a manual “run now” invocation, **When** the user requests it, **Then** the pipeline starts one ingestion run and reports a clear success/failure result.

---

### User Story 3 - Detect missing data (Priority: P3)

As an operator or developer responsible for data freshness, I want runnable examples that (a) poll for data availability (such as a file appearing) and (b) raise an alert when expected data does not arrive by a deadline, so I can implement basic SLA monitoring alongside ingestion.

**Why this priority**: Polling and watchdog behavior are common in batch/ETL-style systems and are frequent sources of subtle reliability issues.

**Independent Test**: Can be fully tested by running the file-polling example with a known target and by running the watchdog example with a short deadline and observing it alerts when no data arrives.

**Acceptance Scenarios**:

1. **Given** a configured file path pattern and polling interval, **When** the file becomes available, **Then** ingestion initiates promptly and does not re-trigger for the same file unless explicitly configured.
2. **Given** a configured deadline, **When** no data-availability event occurs before the deadline, **Then** the example emits an explicit “deadline missed” alert signal.

---

[Add more user stories as needed, each with an assigned priority]

### Edge Cases

- Trigger drift: interval schedules should not accumulate unbounded delay over time.
- Time edge cases: daily trigger behavior should be well-defined across day boundaries and clock changes.
- Duplicate events: webhook or file-availability triggers may deliver duplicates; examples must show a safe default behavior.
- Backpressure: triggers can fire faster than ingestion completes; examples must define whether they queue, coalesce, or skip (default: coalesce to one pending run).
- Partial ingestion: if ingestion starts but the datasource becomes unavailable mid-run, examples must surface a clear failure.
- Large payloads: webhook bodies and files may exceed expected sizes; examples must have a clear limit and error path.

## Requirements *(mandatory)*

## Constitution Compliance *(mandatory)*

Summarize how this feature complies with `.specify/memory/constitution.md`:

- Library-first: standalone reusable package boundary
- Test-first: tests written first and initially failing
- Core independence: if touching `core`, confirm stdlib-only imports
- Quality gates: formatting, linting/static analysis, security checks, tests, coverage
- Docs/examples: at least one runnable example

This feature is example-focused and keeps the library boundary intact. Examples will be runnable, self-contained, and verifiable via automated checks. Where a datasource cannot be accessed without third-party dependencies, examples will demonstrate ingestion via dependency-free mechanisms (for example, reading exported files or using HTTP endpoints) rather than introducing new dependencies.

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: The repository MUST include runnable examples that demonstrate each of the following trigger styles: interval schedule, daily at $HH:MM$, webhook/API event, file-availability polling, manual “run now”, and watchdog/SLA deadline.
- **FR-002**: Each trigger example MUST initiate ingestion by sending a message/event to a “process” step (or equivalent) that performs datasource connection and ingestion.
- **FR-003**: Trigger examples MUST clearly define behavior when a trigger fires while a previous ingestion run is still in progress. The default behavior MUST be **coalesce**: allow at most one pending run; multiple triggers collapse into one follow-up run.
- **FR-004**: The repository MUST include runnable examples that demonstrate ingestion from each of the following datasource categories: files/object storage (via dependency-free access patterns, e.g., local folder stand-in), databases, data warehouses, REST APIs, and observability data (metrics, traces, logs).
- **FR-005**: Examples MUST NOT add third-party libraries; if a datasource category typically requires third-party dependencies, the example MUST use a dependency-free approach (such as ingesting exported data, using vendor-neutral file formats, or consuming an HTTP endpoint).
- **FR-006**: Webhook/API-triggered ingestion MUST validate that an incoming request is well-formed and MUST respond with a clear success/failure outcome.
- **FR-006a**: The webhook example MUST document how to send a local test request and what output indicates a successful cycle.
- **FR-007**: File-availability polling MUST detect new data within a configurable polling interval and MUST avoid re-ingesting the same file by default.
- **FR-008**: Watchdog/SLA behavior MUST produce an explicit alert signal event in the pipeline when expected data does not arrive before a configured deadline. The example MUST still exit 0, but MUST print a clear “ALERT” outcome.
- **FR-009**: Each example MUST produce human-readable output that makes it obvious when ingestion ran, what it ingested (at a high level), and whether it succeeded or failed.
- **FR-010**: Each example MUST be executable in a development environment without additional services by default, using simulated data where needed.
- **FR-011**: Examples MUST be organized as two runnable sets: (a) trigger-focused examples (one per trigger style) using a simple simulated datasource, and (b) datasource-focused examples (one per datasource category) using a simple manual/run-now trigger.

### Assumptions

- Examples default to local, simulated, or self-hosted inputs unless the user explicitly configures external endpoints.
- “Object storage” is demonstrated via a dependency-free local-folder stand-in by default.
- “Databases” and “data warehouses” are demonstrated without in-process drivers; examples may ingest exported extracts (files) or query results exposed via HTTP.
- “Observability data” is demonstrated via dependency-free inputs such as log files, metrics over HTTP, or trace/log exports in a common interchange format.

### Out of Scope

- Implementing vendor-specific authentication/signing flows that would require non-standard dependencies.
- Adding new third-party dependencies to the repository to support specific database, warehouse, or observability protocols.

### Key Entities *(include if feature involves data)*

- **Trigger Configuration**: Defines when/why ingestion should run (interval, daily time, webhook event, polling settings, manual run, deadline policy).
- **Trigger Event**: A message representing a trigger occurrence, including any payload needed to drive ingestion.
- **Ingestion Request**: A normalized request that instructs the process step what data to ingest.
- **Datasource Descriptor**: High-level description of a datasource category and how it is accessed in the example.
- **Ingestion Result**: A success/failure outcome plus basic ingestion summary (counts/timestamps).
- **Alert Signal**: A clear representation that a deadline/SLA was missed.

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: At least 6 runnable trigger examples exist (one per listed trigger style) and each demonstrates a clear successful cycle in under 30 seconds in a default configuration (service-like examples such as webhook must be manually exercisable and show a clear success path).
- **SC-002**: At least 5 runnable datasource examples exist (covering all listed datasource categories), each producing an ingestion summary and a clear success/failure status.
- **SC-003**: Zero third-party libraries are required to build and run all examples.
- **SC-004**: For each example, a user can run it and observe a deterministic outcome (success/failure/alert) without requiring external infrastructure (using simulated or local inputs by default). Exception: the REST API datasource example requires internet access and may fail if the endpoint is unavailable.
