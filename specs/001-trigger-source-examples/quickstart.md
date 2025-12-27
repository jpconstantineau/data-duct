# Quickstart: Trigger & Source Examples

This feature adds runnable examples (Go 1.25, Windows-friendly, stdlib-only) that demonstrate:

- Trigger patterns (interval, daily HH:MM, webhook, file polling, manual run-now, watchdog)
- Datasource ingestion patterns (files/object storage stand-in, database exports, warehouse exports, REST API, observability)

## Prereqs

- Go 1.25+
- Internet access (only required for the REST API datasource example)

## Run trigger-focused examples

(Names below match the planned implementation in `plan.md`.)

- Interval:
  - `go run ./cmd/trigger-interval -interval 1s -duration 10s`
- Daily HH:MM (runs once):
  - `go run ./cmd/trigger-daily -time 23:59 -once`
- Webhook:
  - `go run ./cmd/trigger-webhook -listen 127.0.0.1:8080 -path /trigger`
  - In another terminal:
    - `powershell -Command "Invoke-RestMethod -Method Post -Uri http://127.0.0.1:8080/trigger -ContentType application/json -Body '{\"note\":\"run\"}'"`
- File polling:
  - `go run ./cmd/trigger-file-poll -watch . -poll 1s`
- Manual run-now:
  - `go run ./cmd/trigger-manual-run-now -run-now`
- Watchdog/SLA:
  - `go run ./cmd/trigger-watchdog-sla -deadline 5s -poll 1s`

## Run datasource-focused examples

- Files / object storage stand-in:
  - `go run ./cmd/source-files -root .`
- Database export ingestion:
  - `go run ./cmd/source-database-export -input .\testdata\trigger-source-examples\exports\db-export.csv`
- Warehouse export ingestion:
  - `go run ./cmd/source-warehouse-export -input .\testdata\trigger-source-examples\exports\warehouse-export.csv`
- REST API ingestion (uselessfacts):
  - `go run ./cmd/source-rest-uselessfacts -url "https://uselessfacts.jsph.pl/api/v2/facts/random?language=en"`
- Observability ingestion:
  - `go run ./cmd/source-observability -logs .\testdata\trigger-source-examples\observability\app.log -traces .\testdata\trigger-source-examples\observability\traces.json`

## Expected output (high level)

Each example prints:

- Trigger fired (or manual run invoked)
- Ingestion started
- Ingestion summary (what was ingested)
- Final outcome: `succeeded`, `failed`, or `ALERT`
