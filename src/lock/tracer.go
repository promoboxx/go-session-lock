package lock

import (
	"context"
	"time"
)

// Tracer can trace the flow of calls via child spans.  This is meant to play nice with open tracing
type Tracer interface {
	StartSpanWithContext(ctx context.Context, name string) (Span, context.Context)
	BackgroundRate(sessionID, jobName string, params map[string]string, value int64) error
	BackgroundError(sessionID, jobName string, params map[string]string, code, message string, value int64) error
	BackgroundDuration(sessionID, jobName string, params map[string]string, value time.Duration) error
	BackgroundCustom(sessionID string, jobName string, customName string, params, other map[string]string, value int64) error
}

// Span can hold an error and be finalized.  This is meant to play nice with open tracing
type Span interface {
	Finish()
	SetError(err error)
	SetTag(key string, value interface{})
}

// newNoopTracer exposes a noop tracer that does nothing but fulfill the Tracer interface
func newNoopTracer() Tracer {
	return noopTracer{}
}

type noopTracer struct{}
type noopSpan struct{}

func (noopSpan) SetError(err error)                   {}
func (noopSpan) SetTag(key string, value interface{}) {}

func (noopTracer) StartSpanWithContext(ctx context.Context, name string) (Span, context.Context) {
	return noopSpan{}, ctx
}

func (noopTracer) BackgroundRate(sessionID, jobName string, params map[string]string, value int64) error {
	return nil
}

func (noopTracer) BackgroundError(sessionID, jobName string, params map[string]string, code, message string, value int64) error {
	return nil
}

func (noopTracer) BackgroundDuration(sessionID, jobName string, params map[string]string, value time.Duration) error {
	return nil
}

func (noopTracer) BackgroundCustom(sessionID string, jobName string, customName string, params, other map[string]string, value int64) error {
	return nil
}

func (noopSpan) Finish() {}
