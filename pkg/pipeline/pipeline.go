package pipeline

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jpconstantineau/data-duct/internal/pipelineinternal"
)

type stageKind int

const (
	stageSingle stageKind = iota
	stageBatch
)

type stageDef struct {
	kind        stageKind
	name        string
	buffer      int
	concurrency int

	single      pipelineinternal.SingleHandler
	batch       pipelineinternal.BatchHandler
	batchPolicy pipelineinternal.BatchPolicy
}

type definition struct {
	name   string
	buffer int
	logger pipelineinternal.Logger

	source pipelineinternal.Source
	stages []stageDef
	sink   pipelineinternal.Sink

	currentType reflect.Type
}

// Pipeline is a pipeline builder. Stages can be added in sequence via Then/ThenBatch.
//
// Handlers are provided as ordinary, typed functions and are adapted internally via reflection.
// Example:
//
//	runnable := pipeline.New("example", src).Then(proc).To(sink)
//
// Signature requirements:
//
//	Then:      func(context.Context, In) (Out, error)
//	ThenBatch: func(context.Context, []In) ([]Out, error)
//	To:        func(context.Context, In) error
type Pipeline struct {
	def *definition
}

// New creates a new pipeline builder.
func New[T any](name string, source SourceFunc[T], opts ...Option) *Pipeline {
	o := defaultPipelineOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}

	currentType := typeOf[T]()

	def := &definition{
		name:        name,
		buffer:      o.buffer,
		logger:      pipelineinternal.FromSlog(o.logger),
		currentType: currentType,
		source: func(ctx context.Context) (<-chan any, error) {
			ch, err := source(ctx)
			if err != nil {
				return nil, err
			}

			out := make(chan any)
			go func() {
				defer close(out)
				for {
					select {
					case <-ctx.Done():
						return
					case v, ok := <-ch:
						if !ok {
							return
						}
						select {
						case <-ctx.Done():
							return
						case out <- v:
						}
					}
				}
			}()
			return out, nil
		},
	}

	return &Pipeline{def: def}
}

func (p *Pipeline) Then(handler any, opts ...StageOption) *Pipeline {
	if p == nil || p.def == nil {
		panic("pipeline: builder must not be nil")
	}
	wrapped, outType := wrapSingleHandler(handler, p.def.currentType)

	so := defaultStageOptions()
	so.buffer = p.def.buffer
	for _, opt := range opts {
		if opt != nil {
			opt(&so)
		}
	}

	p.def.stages = append(p.def.stages, stageDef{
		kind:        stageSingle,
		name:        so.name,
		buffer:      so.buffer,
		concurrency: so.concurrency,
		single:      wrapped,
	})

	p.def.currentType = outType
	return p
}

func (p *Pipeline) ThenBatch(handler any, batch BatchPolicy, opts ...StageOption) *Pipeline {
	if p == nil || p.def == nil {
		panic("pipeline: builder must not be nil")
	}

	bp := pipelineinternal.BatchPolicy{Size: batch.Size, MaxWait: batch.MaxWait}
	if bp.Size < 1 {
		panic("pipeline: batch size must be >= 1")
	}

	wrapped, outType := wrapBatchHandler(handler, p.def.currentType)

	so := defaultStageOptions()
	so.buffer = p.def.buffer
	for _, opt := range opts {
		if opt != nil {
			opt(&so)
		}
	}

	p.def.stages = append(p.def.stages, stageDef{
		kind:        stageBatch,
		name:        so.name,
		buffer:      so.buffer,
		concurrency: so.concurrency,
		batch:       wrapped,
		batchPolicy: bp,
	})

	p.def.currentType = outType
	return p
}

func (p *Pipeline) To(sink any, opts ...StageOption) *Runnable {
	if p == nil || p.def == nil {
		panic("pipeline: builder must not be nil")
	}

	wrapped := wrapSink(sink, p.def.currentType)

	so := defaultStageOptions()
	so.buffer = p.def.buffer
	for _, opt := range opts {
		if opt != nil {
			opt(&so)
		}
	}

	p.def.sink = wrapped
	return &Runnable{def: p.def}
}

var (
	ctxType   = reflect.TypeOf((*context.Context)(nil)).Elem()
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

func typeOf[T any]() reflect.Type {
	var zero *T
	return reflect.TypeOf(zero).Elem()
}

func wrapSingleHandler(handler any, expectedIn reflect.Type) (pipelineinternal.SingleHandler, reflect.Type) {
	if handler == nil {
		panic("pipeline: handler must not be nil")
	}
	v := reflect.ValueOf(handler)
	if v.Kind() != reflect.Func {
		panic(fmt.Sprintf("pipeline: handler must be a func, got %T", handler))
	}
	t := v.Type()
	if t.NumIn() != 2 || t.In(0) != ctxType {
		panic(fmt.Sprintf("pipeline: handler must have signature func(context.Context, In) (Out, error), got %s", t.String()))
	}
	if t.NumOut() != 2 || t.Out(1) != errorType {
		panic(fmt.Sprintf("pipeline: handler must have signature func(context.Context, In) (Out, error), got %s", t.String()))
	}
	inType := t.In(1)
	if expectedIn != nil && !expectedIn.AssignableTo(inType) && !expectedIn.ConvertibleTo(inType) {
		panic(fmt.Sprintf("pipeline: handler input type %s is not compatible with previous stage output %s", inType, expectedIn))
	}
	outType := t.Out(0)

	wrapped := func(ctx context.Context, input any) (any, error) {
		inVal, err := adaptValue(input, inType)
		if err != nil {
			return nil, err
		}
		outs := v.Call([]reflect.Value{reflect.ValueOf(ctx), inVal})
		if !outs[1].IsNil() {
			return nil, outs[1].Interface().(error)
		}
		return outs[0].Interface(), nil
	}

	return wrapped, outType
}

func wrapBatchHandler(handler any, expectedElem reflect.Type) (pipelineinternal.BatchHandler, reflect.Type) {
	if handler == nil {
		panic("pipeline: batch handler must not be nil")
	}
	v := reflect.ValueOf(handler)
	if v.Kind() != reflect.Func {
		panic(fmt.Sprintf("pipeline: batch handler must be a func, got %T", handler))
	}
	t := v.Type()
	if t.NumIn() != 2 || t.In(0) != ctxType || t.In(1).Kind() != reflect.Slice {
		panic(fmt.Sprintf("pipeline: batch handler must have signature func(context.Context, []In) ([]Out, error), got %s", t.String()))
	}
	if t.NumOut() != 2 || t.Out(1) != errorType || t.Out(0).Kind() != reflect.Slice {
		panic(fmt.Sprintf("pipeline: batch handler must have signature func(context.Context, []In) ([]Out, error), got %s", t.String()))
	}
	inSliceType := t.In(1)
	inElem := inSliceType.Elem()
	if expectedElem != nil && !expectedElem.AssignableTo(inElem) && !expectedElem.ConvertibleTo(inElem) {
		panic(fmt.Sprintf("pipeline: batch handler input element type %s is not compatible with previous stage output %s", inElem, expectedElem))
	}
	outSliceType := t.Out(0)
	outElem := outSliceType.Elem()

	wrapped := func(ctx context.Context, inputs []any) ([]any, error) {
		slice := reflect.MakeSlice(inSliceType, len(inputs), len(inputs))
		for i := range inputs {
			iv, err := adaptValue(inputs[i], inElem)
			if err != nil {
				return nil, err
			}
			slice.Index(i).Set(iv)
		}
		outs := v.Call([]reflect.Value{reflect.ValueOf(ctx), slice})
		if !outs[1].IsNil() {
			return nil, outs[1].Interface().(error)
		}
		outSlice := outs[0]
		anyOut := make([]any, 0, outSlice.Len())
		for i := 0; i < outSlice.Len(); i++ {
			ov := outSlice.Index(i)
			// Ensure interface materialization preserves concrete type.
			if ov.Type().AssignableTo(outElem) {
				anyOut = append(anyOut, ov.Interface())
			} else {
				anyOut = append(anyOut, ov.Convert(outElem).Interface())
			}
		}
		return anyOut, nil
	}

	return wrapped, outElem
}

func wrapSink(sink any, expectedIn reflect.Type) pipelineinternal.Sink {
	if sink == nil {
		panic("pipeline: sink must not be nil")
	}
	v := reflect.ValueOf(sink)
	if v.Kind() != reflect.Func {
		panic(fmt.Sprintf("pipeline: sink must be a func, got %T", sink))
	}
	t := v.Type()
	if t.NumIn() != 2 || t.In(0) != ctxType {
		panic(fmt.Sprintf("pipeline: sink must have signature func(context.Context, In) error, got %s", t.String()))
	}
	if t.NumOut() != 1 || t.Out(0) != errorType {
		panic(fmt.Sprintf("pipeline: sink must have signature func(context.Context, In) error, got %s", t.String()))
	}
	inType := t.In(1)
	if expectedIn != nil && !expectedIn.AssignableTo(inType) && !expectedIn.ConvertibleTo(inType) {
		panic(fmt.Sprintf("pipeline: sink input type %s is not compatible with previous stage output %s", inType, expectedIn))
	}

	return func(ctx context.Context, input any) error {
		inVal, err := adaptValue(input, inType)
		if err != nil {
			return err
		}
		outs := v.Call([]reflect.Value{reflect.ValueOf(ctx), inVal})
		if outs[0].IsNil() {
			return nil
		}
		return outs[0].Interface().(error)
	}
}

func adaptValue(input any, target reflect.Type) (reflect.Value, error) {
	if input == nil {
		// nil can only be passed to some types; for others, reject.
		switch target.Kind() {
		case reflect.Interface, reflect.Pointer, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
			return reflect.Zero(target), nil
		default:
			return reflect.Value{}, fmt.Errorf("pipeline: cannot pass nil to %s", target)
		}
	}

	v := reflect.ValueOf(input)
	if v.Type().AssignableTo(target) {
		return v, nil
	}
	if v.Type().ConvertibleTo(target) {
		return v.Convert(target), nil
	}
	return reflect.Value{}, fmt.Errorf("pipeline: cannot use %T as %s", input, target)
}
