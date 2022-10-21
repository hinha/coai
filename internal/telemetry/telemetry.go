package telemetry

import (
	"context"
	"github.com/hinha/coai/config"
	"github.com/hinha/coai/internal/logger"
	"github.com/hinha/coai/internal/telemetry/exporter"
	"github.com/hinha/coai/internal/telemetry/metric"
	"github.com/hinha/coai/internal/telemetry/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

var Module = fx.Module("open_telemetry",
	fx.Provide(func(config *config.Config) Config {
		return Config{
			Endpoint:    config.Otel.Server,
			ServiceName: config.Server.Name,
			Enable:      config.Otel.Enabled,
		}
	}),
	fx.Provide(NewTelemetry),
)

type Config struct {
	Endpoint    string
	ServiceName string
	Enable      bool
}

type OpenTelemetry struct {
	closers []func(ctx context.Context) error
}

func NewTelemetry(config Config, logger *logger.Logger) *OpenTelemetry {
	tm := &OpenTelemetry{}

	if !config.Enable {
		return tm
	}

	ctx := context.Background()
	metricClient := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(config.Endpoint))

	metricExporter, err := otlpmetric.New(ctx, metricClient)
	if err != nil {
		logger.Console().Fatal("Failed to create the collector metric exporter", zap.Error(err))
	}
	logger.Console().Debug("Finish create the collector metric exporter")

	pusher, pusherCloseFn, err := metric.NewMeterProviderBuilder(logger).
		SetExporter(metricExporter).
		SetHistogramBoundaries([]float64{5, 10, 25, 50, 100, 200, 400, 800, 1000}).
		Build()
	if err != nil {
		logger.Console().Fatal("failed initializing the metric provider", zap.Error(err))
	}
	tm.closers = append(tm.closers, pusherCloseFn)
	global.SetMeterProvider(pusher)

	spanExporter := exporter.NewTraceOTLP(config.Endpoint, logger)
	tracerProvider, tracerProviderCloseFn, err := trace.NewTraceProviderBuilder(config.ServiceName, logger).
		SetExporter(spanExporter).
		Build()
	if err != nil {
		logger.Console().Fatal("failed initializing the tracer provider", zap.Error(err))
	}
	tm.closers = append(tm.closers, tracerProviderCloseFn)

	// set global propagator to trace context (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return tm
}

func (t *OpenTelemetry) Close() error {
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		cancel()
	}()

	for _, f := range t.closers {
		if err := f(ctxShutDown); err != nil {
			return err
		}
	}

	return nil
}
