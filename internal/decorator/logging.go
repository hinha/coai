package decorator

import (
	"context"
	"github.com/hinha/coai/internal/logger"
	"go.uber.org/zap"
)

type commandLoggingDecorator[C any] struct {
	base   CommandHandler[C]
	logger *logger.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	log := d.logger.Core().With(
		zap.String("command", handlerType),
		zap.Reflect("command_body", cmd),
	)

	defer func() {
		log.DebugCtx(ctx, "Command executed", zap.Error(err))
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[C any, R any] struct {
	base   QueryHandler[C, R]
	logger *logger.Logger
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, query C) (result R, err error) {
	log := d.logger.Core().With(
		zap.String("query", generateActionName(query)),
		zap.Reflect("query_body", query),
	)

	defer func() {
		log.DebugCtx(ctx, "Query executed", zap.Error(err))
	}()

	return d.base.Handle(ctx, query)
}
