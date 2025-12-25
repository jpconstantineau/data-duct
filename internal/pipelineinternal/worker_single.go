package pipelineinternal

import (
	"context"
	"sync"
)

func workerSingle(ctx context.Context, in <-chan feed, out chan<- feed, handler SingleHandler, concurrency int, logger Logger) {
	if concurrency < 1 {
		concurrency = 1
	}

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for f := range in {
				// Respect cancellation.
				select {
				case <-ctx.Done():
					// Drain input by continuing the range; but stop processing.
					continue
				default:
				}

				outData, err := handler(ctx, f.Data)
				if err != nil {
					// Do not emit an output item for this failed input.
					continue
				}

				nf := feed{RootCtx: f.RootCtx, PipelineName: f.PipelineName, Data: outData}
				select {
				case <-ctx.Done():
					continue
				case out <- nf:
				}
			}
		}()
	}

	wg.Wait()
	close(out)
	logger.Debug("pipeline stage complete")
}
