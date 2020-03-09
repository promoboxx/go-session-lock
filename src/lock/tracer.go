package lock

import (
	"context"
)

// Tracer can trace the flow of calls via child spans.  This is meant to play nice with open tracing
type Tracer interface {
	StartSpanWithContext(ctx context.Context, name string) (Span, context.Context)
}

// Span can hold an error and be finalized.  This is meant to play nice with open tracing
type Span interface {
	Finish()
}

// newNoopTracer exposes a noop tracer that does nothing but fulfill the Tracer interface
func newNoopTracer() Tracer {
	return noopTracer{}
}

type noopTracer struct{}
type noopSpan struct{}

func (noopTracer) StartSpanWithContext(ctx context.Context, name string) (Span, context.Context) {
	return noopSpan{}, ctx
}

func (noopSpan) Finish() {}
