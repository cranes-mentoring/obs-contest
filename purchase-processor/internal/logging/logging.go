package logging

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func SetupLogger() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer Logger.Sync()
}

// AddTraceContextToLogger attaches trace context from the provided context to the logger and returns the updated logger.
func AddTraceContextToLogger(ctx context.Context) *zap.Logger {
	spanContext := trace.SpanFromContext(ctx).SpanContext()
	if !spanContext.IsValid() {

		tracer := otel.GetTracerProvider().Tracer("purchase-processor")

		ctx, span := tracer.Start(ctx, "new-span")
		defer span.End()

		spanContext = trace.SpanFromContext(ctx).SpanContext()

		return Logger.With(
			zap.String("trace_id", spanContext.TraceID().String()),
			zap.String("span_id", spanContext.SpanID().String()),
		)
	}

	return Logger.With(
		zap.String("trace_id", spanContext.TraceID().String()),
		zap.String("span_id", spanContext.SpanID().String()),
	)
}
