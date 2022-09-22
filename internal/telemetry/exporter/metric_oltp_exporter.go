package exporter

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.uber.org/zap"

	"github.com/hinha/coai/internal/logger"
)

func NewMetricOTLP(endpoint string, log *logger.Logger) *otlpmetric.Exporter {
	ctx := context.Background()
	metricClient := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endpoint))

	metricExp, err := otlpmetric.New(ctx, metricClient)
	if err != nil {
		log.Console().Fatal("Failed to create the collector metric exporter", zap.Error(err))
	}
	log.Console().Debug("Finish create the collector metric exporter")

	return metricExp
}
