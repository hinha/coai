package decorator

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/hinha/coai/internal/logger"
)

type commandLoggingDecorator[C any] struct {
	base   CommandHandler[C]
	logger *logger.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	log := d.logger.Core().With(
		zap.String("command", handlerType),
		zap.String("command_body", fmt.Sprintf("%#v", cmd)),
	)

	log.Debug("Executing command")
	defer func() {
		if err == nil {
			log.InfoCtx(ctx, "Command executed successfully")
		} else {
			log.ErrorCtx(ctx, "Failed to execute command", zap.Error(err))
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[C any, R any] struct {
	base   QueryHandler[C, R]
	logger *logger.Logger
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	log := d.logger.Core().With(
		zap.String("query", generateActionName(cmd)),
		zap.String("query_body", fmt.Sprintf("%#v", cmd)),
	)

	log.Debug("Executing query")
	defer func() {
		if err == nil {
			log.InfoCtx(ctx, "Query executed successfully")
		} else {
			log.ErrorCtx(ctx, "Failed to execute query")
		}
	}()

	return d.base.Handle(ctx, cmd)
}
