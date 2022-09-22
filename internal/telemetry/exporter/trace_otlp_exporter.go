package exporter

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/hinha/coai/internal/logger"
)

func NewTraceOTLP(endpoint string, log *logger.Logger) *otlptrace.Exporter {
	ctx := context.Background()
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))

	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		log.Console().Fatal("Failed to create the collector trace exporter", zap.Error(err))
	}
	log.Console().Debug("Finish create the collector trace exporter")

	return traceExp
}
