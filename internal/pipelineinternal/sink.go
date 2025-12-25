package pipelineinternal

import "context"

func sinkConsume(ctx context.Context, in <-chan feed, sink Sink, policy *errorPolicy, logger Logger) {
	for f := range in {
		// Always drain to avoid blocking upstream, even after failure.
		if ctx.Err() != nil {
			continue
		}
		if policy.get() != nil {
			continue
		}
		if err := sink(ctx, f.Data); err != nil {
			policy.set(err)
			logger.Error("pipeline sink error", "error", err)
		}
	}
}
