package pipelineinternal

import (
	"errors"
	"sync"
)

var ErrInvalidConfig = errors.New("pipeline: invalid configuration")

// errorPolicy captures the first failure cause.
type errorPolicy struct {
	once  sync.Once
	cause error
}

func (p *errorPolicy) set(err error) {
	if err == nil {
		return
	}
	p.once.Do(func() { p.cause = err })
}

func (p *errorPolicy) get() error {
	return p.cause
}
