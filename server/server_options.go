package server

import (
	"github.com/hinha/coai/internal/telemetry/exporter"
	"github.com/hinha/coai/internal/telemetry/metric"
	"github.com/hinha/coai/internal/telemetry/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

type Option func(s *Server)

func WithTelemetry(name, endpoint string) Option {
	return func(s *Server) {
		if s.logger == nil {
			panic("logger not set")
		}

		if endpoint == "" {
			return
		}

		metricExporter := exporter.NewMetricOTLP(endpoint, s.logger)
		pusher, pusherCloseFn, err := metric.NewMeterProviderBuilder().
			SetExporter(metricExporter).
			SetHistogramBoundaries([]float64{5, 10, 25, 50, 100, 200, 400, 800, 1000}).
			Build()
		if err != nil {
			s.logger.LogDefault().Error("failed initializing the metric provider", zap.Error(err))
		}
		s.Closers = append(s.Closers, pusherCloseFn)
		global.SetMeterProvider(pusher)

		spanExporter := exporter.NewTraceOTLP(endpoint, s.logger)
		tracerProvider, tracerProviderCloseFn, err := trace.NewTraceProviderBuilder(name, s.logger).
			SetExporter(spanExporter).
			Build()
		if err != nil {
			s.logger.LogDefault().Error("failed initializing the tracer provider", zap.Error(err))
		}
		s.Closers = append(s.Closers, tracerProviderCloseFn)

		// set global propagator to tracecontext (the default is no-op).
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		otel.SetTracerProvider(tracerProvider)
	}
}
