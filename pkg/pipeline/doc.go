// Package pipeline provides a small, typed pipeline builder and runner.
//
// The public surface is intentionally minimal:
//
//   - New: define pipeline name and source
//   - Then / ThenBatch: add processors
//   - To: attach a sink
//   - Run: execute with a root context
package pipeline
