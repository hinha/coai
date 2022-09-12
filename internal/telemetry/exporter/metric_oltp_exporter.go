package exporter

import (
	"context"
	"github.com/hinha/coai/logger"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.uber.org/zap"
)

func NewMetricOTLP(endpoint string, log *logger.Logger) *otlpmetric.Exporter {
	ctx := context.Background()
	metricClient := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endpoint))

	metricExp, err := otlpmetric.New(ctx, metricClient)
	if err != nil {
		log.LogDefault().Error("Failed to create the collector metric exporter", zap.Error(err))
	}

	return metricExp
}
