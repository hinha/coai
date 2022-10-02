package query

import (
	"context"
	"github.com/hinha/coai/internal/decorator"
	"github.com/hinha/coai/internal/logger"
)

type AllUsers struct{}

type AllUserHandler decorator.QueryHandler[AllUsers, []User]

type allTrainingsHandler struct {
	readModel AllUsersReadModel
}

func NewAllUserHandler(
	readModel AllUsersReadModel,
	logger *logger.Logger,
	metricsClient decorator.MetricsClient,
) decorator.QueryHandler[AllUsers, []User] {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[AllUsers, []User](
		allTrainingsHandler{readModel: readModel},
		logger,
		metricsClient,
	)
}

func (l allTrainingsHandler) Handle(ctx context.Context, _ AllUsers) (tr []User, err error) {
	return l.readModel.AllUsers(ctx)
}

type AllUsersReadModel interface {
	AllUsers(ctx context.Context) ([]User, error)
}
