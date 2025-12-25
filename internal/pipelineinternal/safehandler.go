package pipelineinternal

import (
	"context"
	"fmt"
)

func safeSingle(name string, h SingleHandler, policy *errorPolicy) SingleHandler {
	return func(ctx context.Context, input any) (out any, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("pipeline: panic in handler%s: %v", formatStage(name), r)
				policy.set(err)
			}
		}()

		out, err = h(ctx, input)
		if err != nil {
			policy.set(err)
		}
		return out, err
	}
}

func safeBatch(name string, h BatchHandler, policy *errorPolicy) BatchHandler {
	return func(ctx context.Context, inputs []any) (outs []any, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("pipeline: panic in batch handler%s: %v", formatStage(name), r)
				policy.set(err)
			}
		}()

		outs, err = h(ctx, inputs)
		if err != nil {
			policy.set(err)
		}
		return outs, err
	}
}

func formatStage(name string) string {
	if name == "" {
		return ""
	}
	return " (" + name + ")"
}
