package command

import (
	"context"
	"github.com/hinha/coai/core/users/domain/connection"
	"github.com/hinha/coai/internal/decorator"
	"github.com/hinha/coai/internal/logger"
)

type Ping struct {
	StateUP   int
	StateDOWN int
}

type PingHandler decorator.CommandHandler[Ping]

type pingHandlerHandler struct {
	repo connection.Repository
}

func NewPingHandler(
	repo connection.Repository,
	logger *logger.Logger,
	metricsClient decorator.MetricsClient,
) PingHandler {
	if repo == nil {
		panic("nil userGroup repo")
	}

	return decorator.ApplyCommandDecorators[Ping](
		pingHandlerHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (p pingHandlerHandler) Handle(ctx context.Context, cmd Ping) error {
	if err := p.repo.Ping(); err != nil {
		// TODO send alert
		return err
	}
	return nil
}
