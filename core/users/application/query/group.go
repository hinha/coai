package query

import (
	"context"
	"github.com/hinha/coai/core/users/domain"
	"github.com/hinha/coai/core/users/domain/user_groups"
	"github.com/hinha/coai/internal/decorator"
	"github.com/hinha/coai/internal/logger"
)

type GetGroup struct {
	Name string
	Id   int64
}

type GetGroupHandler decorator.QueryHandler[GetGroup, *domain.UserGroup]

type getGroupHandler struct {
	userGroups user_groups.Repository
}

func NewGetGroupHandler(
	userGroups user_groups.Repository,
	logger *logger.Logger,
	metricsClient decorator.MetricsClient,
) GetGroupHandler {
	if userGroups == nil {
		panic("nil userGroups")
	}

	return decorator.ApplyQueryDecorators[GetGroup, *domain.UserGroup](
		getGroupHandler{userGroups: userGroups},
		logger,
		metricsClient,
	)
}

func (h getGroupHandler) Handle(ctx context.Context, query GetGroup) (*domain.UserGroup, error) {
	group, err := h.userGroups.GetUserGroup(ctx, query.Id, query.Name)
	if err != nil {
		return nil, err
	}

	return group, nil
}
