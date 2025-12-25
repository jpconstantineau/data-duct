package pipelineinternal

// Cancellation behavior is handled by consistently selecting on ctx.Done() in
// all pumps/workers and by ensuring Run waits for all goroutines to exit.
