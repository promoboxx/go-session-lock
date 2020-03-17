package lock

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// Tracer can trace the flow of calls via child spans.  This is meant to play nice with open tracing
type Tracer interface {
	StartSpanWithContext(ctx context.Context, name string) (opentracing.Span, context.Context)
}

// newNoopTracer exposes a noop tracer that does nothing but fulfill the Tracer interface
func newNoopTracer() Tracer {
	return noopTracer{}
}

type noopTracer struct{}

func (noopTracer) StartSpanWithContext(ctx context.Context, name string) (opentracing.Span, context.Context) {
	span := opentracing.SpanFromContext(ctx)
	return span, opentracing.ContextWithSpan(ctx, span)
}
