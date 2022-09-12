package exporter

import (
	"context"
	"github.com/hinha/coai/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func NewTraceOTLP(endpoint string, log *logger.Logger) *otlptrace.Exporter {
	ctx := context.Background()
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))

	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		log.LogDefault().Error("Failed to create the collector trace exporter", zap.Error(err))
	}

	return traceExp
}
