package trace

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
)

type ProviderConfig struct {
	JaegerEndpoint string
	ServiceName    string
	ServiceVersion string
	Environment    string
	Disabled       bool
}

type Provider struct {
	provider trace.TracerProvider
}

func NewStdOutExporter(w io.Writer) (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func NewJaegerExporter(jaegerEndpoint string) (sdktrace.SpanExporter, error) {
	return jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)),
	)
}

func NewResource(config ProviderConfig) *sdkresource.Resource {
	return sdkresource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.ServiceVersionKey.String(config.ServiceVersion),
		semconv.DeploymentEnvironmentKey.String(config.Environment),
	)
}

func NewProvider(config ProviderConfig, exporter sdktrace.SpanExporter, resource *sdkresource.Resource) (Provider, error) {
	if config.Disabled {
		return Provider{provider: trace.NewNoopTracerProvider()}, nil
	}
	prv := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(prv)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return Provider{provider: prv}, nil
}

func (p Provider) Close(ctx context.Context) {
	if prv, ok := p.provider.(*sdktrace.TracerProvider); ok {
		if err := prv.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}
}
