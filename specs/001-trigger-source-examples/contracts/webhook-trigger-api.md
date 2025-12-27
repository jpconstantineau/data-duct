# Contract: Webhook Trigger API (Example)

**Date**: 2025-12-26

This contract describes the local HTTP API used by the webhook-trigger runnable example.

## Endpoint

- **Method**: `POST`
- **Path**: `/trigger`
- **Content-Type**: `application/json`

## Request

### Schema

```json
{
  "event_id": "string (optional)",
  "source_category": "string (optional)",
  "note": "string (optional)"
}
```

### Validation

- If present, `event_id` must be a non-empty string.
- If present, `source_category` must be one of:
  - `files`
  - `object_storage`
  - `database`
  - `warehouse`
  - `rest_api`
  - `observability`

## Response

### Success

- **Status**: `202 Accepted`

```json
{
  "accepted": true,
  "event_id": "string",
  "message": "string"
}
```

### Failure

- **Status**: `400 Bad Request` for validation errors
- **Status**: `500 Internal Server Error` for unexpected errors

```json
{
  "accepted": false,
  "error": "string"
}
```

## Notes

- This API is intended for local development only.
- Implementation must use Go standard library only (`net/http`, `encoding/json`).
