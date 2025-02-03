package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer(ctx context.Context) func() {
	exp, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint("otel-collector:4317"),
			otlptracegrpc.WithInsecure(),
		),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP trace exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("purchase-processor"),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String("stg"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	bsp := trace.NewBatchSpanProcessor(exp)
	tp := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(exp),
		trace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tp)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Printf("OpenTelemetry error: %v", err)
	}))

	return func() {
		err := tp.Shutdown(ctx)
		if err != nil {
			log.Printf("error shutting down tracer provider: %v", err)
		}
	}
}
