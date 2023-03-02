package trace

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// ShutdownProvider attempts to gracefully shutdown a created provider.
// Wait for a max of 5 seconds to prevent the service from hanging and panics on error.
func ShutdownProvider(ctx context.Context, provider *tracesdk.TracerProvider) {
	if provider != nil {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
}

// JaegerProvider creates a new batch Jaeger based TracerProvider exporting to the passed url.
// Errors can be safely ignored as a TraceProvider with no exporter will always be returned in the event of an error.
func JaegerProvider(url string, service string, attr ...attribute.KeyValue) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return NoopProvider(), fmt.Errorf("failed to initialize Jaeger trace exporter: %w", err)
	}
	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			append([]attribute.KeyValue{
				semconv.ServiceName(service),
			}, attr...)...,
		)),
	), nil
}

// LoggingProvider creates a new batch Logging based TracerProvider.
// Errors can be safely ignored as a NoopProvider will always be returned in the event of an error.
func LoggingProvider() (*tracesdk.TracerProvider, error) {
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return NoopProvider(), fmt.Errorf("failed to initialize Logging trace exporter: %w", err)
	}
	return tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithSpanProcessor(tracesdk.NewBatchSpanProcessor(exp)),
	), nil
}

// NoopProvider returns an TraceProvider with no exporters.
func NoopProvider() *tracesdk.TracerProvider {
	return tracesdk.NewTracerProvider()
}
