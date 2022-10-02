package command

import (
	"context"
	"github.com/hinha/coai/core/users/domain"
	"github.com/hinha/coai/internal/decorator"
	"github.com/hinha/coai/internal/logger"
)

type RegisterUser struct {
	User domain.User
}

type RegisterUserHandler decorator.CommandHandler[RegisterUser]

type registerUserHandler struct {
	repo domain.Repository
}

func NewRegisterUserHandler(
	repo domain.Repository,
	logger *logger.Logger,
	metricsClient decorator.MetricsClient) decorator.CommandHandler[RegisterUser] {

	return decorator.ApplyCommandDecorators[RegisterUser](
		registerUserHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h registerUserHandler) Handle(ctx context.Context, cmd RegisterUser) error {
	return h.repo.AddUser(ctx, &cmd.User)
}
