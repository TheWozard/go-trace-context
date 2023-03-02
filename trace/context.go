package trace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type contextKey int

const providerKey contextKey = iota

// NewContext attaches a passed trace.TraceProvider to the context.
func NewContext(ctx context.Context, provider trace.TracerProvider) context.Context {
	return context.WithValue(ctx, providerKey, provider)
}

// From safely returns a trace.TraceProvider from the context.
// If non exits then a trace.NewNoopTracerProvider() is returned instead.
func From(ctx context.Context) trace.TracerProvider {
	if provider, ok := ctx.Value(providerKey).(trace.TracerProvider); ok {
		return provider
	}
	return trace.NewNoopTracerProvider()
}
