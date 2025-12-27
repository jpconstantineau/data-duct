# Contract: REST Datasource (Useless Facts)

**Date**: 2025-12-26

This contract describes the external REST endpoint used by the REST datasource runnable example.

## Endpoint

- **Method**: `GET`
- **URL**: `https://uselessfacts.jsph.pl/api/v2/facts/random?language=en`

## Expected Response

- **Status**: `200 OK`
- **Content-Type**: JSON

### Schema (subset)

```json
{
  "id": "string",
  "text": "string",
  "source": "string",
  "source_url": "string",
  "language": "string",
  "permalink": "string"
}
```

## Error Handling Requirements (example behavior)

- Non-2xx responses must be treated as ingestion failure with a clear message.
- JSON parse errors must be treated as ingestion failure with a clear message.
- Missing/empty `text` must be treated as ingestion failure with a clear message.
- The example must set reasonable HTTP timeouts to avoid hanging.

## Dependencies

- Must be implemented using Go standard library only (`net/http`, `encoding/json`).
