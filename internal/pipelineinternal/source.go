package pipelineinternal

import "context"

func sourcePump(rootCtx context.Context, sourceCtx context.Context, src <-chan any, out chan<- feed, pipelineName string) {
	for {
		select {
		case <-sourceCtx.Done():
			return
		case v, ok := <-src:
			if !ok {
				return
			}
			f := feed{RootCtx: rootCtx, PipelineName: pipelineName, Data: v}
			select {
			case <-sourceCtx.Done():
				return
			case out <- f:
			}
		}
	}
}
