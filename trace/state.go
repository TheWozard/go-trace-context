package trace

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

// State is a json/yaml marshal-able form of the trace.SpanContext
type State struct {
	TraceID string `json:"trace_id" yaml:"trace_id"`
	SpanID  string `json:"span_id,omitempty" yaml:"span_id,omitempty"`
}

// Root returns the root of the current State.
func (s *State) Root() *State {
	if s == nil {
		return nil
	}
	return &State{
		TraceID: s.TraceID,
	}
}

// ContextFromState restores a context from a captured json/yaml marshal-able State.
func ContextFromState(ctx context.Context, state *State) (context.Context, error) {
	if state == nil {
		return ctx, nil
	}
	var cfg trace.SpanContextConfig
	var err error

	if state.TraceID != "" {
		cfg.TraceID, err = trace.TraceIDFromHex(state.TraceID)
		if err != nil {
			return ctx, fmt.Errorf("failed to parse TraceID, %v", err)
		}
	}
	if state.SpanID != "" {
		cfg.SpanID, err = trace.SpanIDFromHex(state.SpanID)
		if err != nil {
			return ctx, fmt.Errorf("failed to parse SpanID, %v", err)
		}
	}

	return trace.ContextWithSpanContext(ctx, trace.NewSpanContext(cfg)), nil
}

// ContextFromState creates a json/yaml marshal-able State from the current context.
// Returns nil if the context contains no Span.
func ContextToState(ctx context.Context) *State {
	span := trace.SpanContextFromContext(ctx)
	if span.TraceID().IsValid() {
		return &State{
			TraceID: span.TraceID().String(),
			SpanID:  span.SpanID().String(),
		}
	}
	return nil
}
