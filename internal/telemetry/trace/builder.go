package trace

import (
	"context"
	"github.com/hinha/coai/logger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.uber.org/zap"
	"time"
)

type CloseFunc func(ctx context.Context) error

func NewTraceProviderBuilder(name string, log *logger.Logger) *traceProviderBuilder {
	return &traceProviderBuilder{
		name: name,
		log:  log,
	}
}

type traceProviderBuilder struct {
	name     string
	exporter trace.SpanExporter
	log      *logger.Logger
}

func (b *traceProviderBuilder) SetExporter(exp trace.SpanExporter) *traceProviderBuilder {
	b.exporter = exp
	return b
}

func (b *traceProviderBuilder) Build() (*trace.TracerProvider, CloseFunc, error) {
	b.log.LogDefault().Info("Tracer provider Builder")

	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		// resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backend
			semconv.ServiceNameKey.String(b.name),
		),
	)
	if err != nil {
		b.log.LogDefault().Error("Something tracer provider Builder", zap.Error(err))
		return nil, nil, err
	}

	bsp := trace.NewBatchSpanProcessor(b.exporter)
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)
	return tracerProvider, func(ctx context.Context) error {
		cxt, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := b.exporter.Shutdown(cxt); err != nil {
			return err
		}
		return err
	}, nil
}
