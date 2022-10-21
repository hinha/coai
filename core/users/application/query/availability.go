package query

import (
	"context"
	"github.com/hinha/coai/core/users/domain/user_groups"
	"github.com/hinha/coai/internal/decorator"
	"github.com/hinha/coai/internal/logger"
)

type GroupAvailability struct {
	Name string
}

type GroupAvailabilityHandler decorator.QueryHandler[GroupAvailability, bool]

type groupAvailabilityHandler struct {
	userGroups user_groups.Repository
}

func NewGroupAvailabilityHandler(
	userGroups user_groups.Repository,
	logger *logger.Logger,
	metricsClient decorator.MetricsClient,
) GroupAvailabilityHandler {
	if userGroups == nil {
		panic("nil userGroups")
	}

	return decorator.ApplyQueryDecorators[GroupAvailability, bool](
		groupAvailabilityHandler{userGroups: userGroups},
		logger,
		metricsClient,
	)
}

func (h groupAvailabilityHandler) Handle(ctx context.Context, query GroupAvailability) (bool, error) {
	group, err := h.userGroups.IsActivateGroup(ctx, query.Name)
	if err != nil {
		return false, err
	}

	return group.IsAvailable(), nil
}
