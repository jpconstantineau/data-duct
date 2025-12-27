# Data Model: Trigger & Source Examples

**Date**: 2025-12-26

This feature is example-focused; entities below describe the data flowing through examples (not a persisted schema).

## Entities

### TriggerConfiguration
Represents how a runnable example decides *when* to initiate ingestion.

- `Kind`: one of `interval`, `daily_hhmm`, `webhook`, `file_poll`, `manual`, `watchdog`
- `Interval`: duration for interval schedules
- `DailyTime`: local time string `HH:MM`
- `Webhook`: host/port/path + validation rules
- `FilePoll`: directory/path pattern + poll interval
- `Deadline`: duration or absolute timestamp by which data must arrive
- `BackpressurePolicy`: default = `coalesce_one_pending`

### TriggerEvent
Represents a single trigger occurrence.

- `EventID`: unique identifier for the event (string)
- `At`: trigger timestamp
- `Kind`: matches `TriggerConfiguration.Kind`
- `Payload`: trigger-specific attributes (e.g., webhook body fields, file path)

### IngestionRequest
A normalized instruction that the “process” step can use to connect to a datasource and ingest data.

- `RequestID`
- `RequestedAt`
- `SourceCategory`: `files`, `object_storage`, `database`, `warehouse`, `rest_api`, `observability`
- `Location`: datasource locator (path/URL) depending on category

### IngestionResult
Outcome of an ingestion run.

- `RequestID`
- `StartedAt` / `FinishedAt`
- `Status`: `succeeded`, `failed`, `alert`
- `ItemsIngested`: count (when meaningful)
- `Summary`: short human-readable summary

### AlertSignal
Represents watchdog/SLA signaling.

- `At`
- `Deadline`
- `ObservedLastArrivalAt`: optional
- `Message`: required

## REST API Contract (Useless Facts)

The REST datasource example fetches a random fact from:

- URL: `https://uselessfacts.jsph.pl/api/v2/facts/random?language=en`
- Method: `GET`

Response (JSON) fields used by the example:

- `id` (string)
- `text` (string)
- `source` (string)
- `source_url` (string)
- `language` (string)
- `permalink` (string)

The example should treat unknown/extra fields as acceptable and should fail clearly if `text` is missing or not a string.
