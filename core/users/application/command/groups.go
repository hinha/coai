package command

import (
	"context"
	"database/sql"
	"github.com/hinha/coai/core/users/domain"
	"github.com/hinha/coai/core/users/domain/user_groups"
	"github.com/hinha/coai/internal/decorator"
	"github.com/hinha/coai/internal/logger"
	"gorm.io/gorm"
	"time"
)

type CreateUserGroup struct {
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	Timestamp time.Time `json:"timestamp"`
}

type CreateUserGroupHandler decorator.CommandHandler[CreateUserGroup]

type createUserGroupHandler struct {
	repo user_groups.Repository
}

func NewCreateUserGroupHandler(
	repo user_groups.Repository,
	logger *logger.Logger,
	metricsClient decorator.MetricsClient,
) CreateUserGroupHandler {
	if repo == nil {
		panic("nil userGroup repo")
	}

	return decorator.ApplyCommandDecorators[CreateUserGroup](
		createUserGroupHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h createUserGroupHandler) Handle(ctx context.Context, cmd CreateUserGroup) error {
	err := h.repo.AddUserGroup(ctx, &domain.UserGroup{
		Name:      cmd.Name,
		Active:    cmd.Active,
		CreatedAt: cmd.Timestamp,
	}, func(group *domain.UserGroup) (*domain.UserGroup, error) {
		return h.repo.GetUserGroup(ctx, 0, group.Name)
	})

	if err != nil {
		return err
	}

	return nil
}

type UpdateUserGroup struct {
	Id        int64     `json:"id"`
	Active    bool      `json:"active"`
	Timestamp time.Time `json:"timestamp"`
}

type UpdateUserGroupHandler decorator.CommandHandler[UpdateUserGroup]

type updateUserGroupHandler struct {
	repo user_groups.Repository
}

func NewUpdateUserGroupHandler(
	repo user_groups.Repository,
	logger *logger.Logger,
	metricsClient decorator.MetricsClient,
) UpdateUserGroupHandler {
	if repo == nil {
		panic("nil userGroup repo")
	}

	return decorator.ApplyCommandDecorators[UpdateUserGroup](
		updateUserGroupHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h updateUserGroupHandler) Handle(ctx context.Context, cmd UpdateUserGroup) error {
	err := h.repo.UpdateUserGroup(ctx, &domain.UserGroup{
		ID:        cmd.Id,
		Active:    cmd.Active,
		UpdatedAt: cmd.Timestamp,
	})

	if err != nil {
		return err
	}

	return nil
}

type DeleteUserGroup struct {
	Id        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

type DeleteUserGroupHandler decorator.CommandHandler[DeleteUserGroup]

type deleteUserGroupHandler struct {
	repo user_groups.Repository
}

func NewDeleteUserGroupHandler(
	repo user_groups.Repository,
	logger *logger.Logger,
	metricsClient decorator.MetricsClient,
) DeleteUserGroupHandler {
	if repo == nil {
		panic("nil userGroup repo")
	}

	return decorator.ApplyCommandDecorators[DeleteUserGroup](
		deleteUserGroupHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h deleteUserGroupHandler) Handle(ctx context.Context, cmd DeleteUserGroup) error {
	err := h.repo.DeleteUserGroup(ctx, &domain.UserGroup{
		ID:        cmd.Id,
		UpdatedAt: cmd.Timestamp,
		DeletedAt: gorm.DeletedAt(sql.NullTime{
			Time:  cmd.Timestamp,
			Valid: true,
		}),
	})

	if err != nil {
		return err
	}

	return nil
}

type ActivateUserGroup struct {
	Id        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

type ActivateUserGroupHandler decorator.CommandHandler[ActivateUserGroup]

type activateUserGroupHandler struct {
	repo user_groups.Repository
}

func NewActivateUserGroupHandler(
	repo user_groups.Repository,
	logger *logger.Logger,
	metricsClient decorator.MetricsClient,
) ActivateUserGroupHandler {
	if repo == nil {
		panic("nil userGroup repo")
	}

	return decorator.ApplyCommandDecorators[ActivateUserGroup](
		activateUserGroupHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h activateUserGroupHandler) Handle(ctx context.Context, cmd ActivateUserGroup) error {
	err := h.repo.UpdateGroupActive(ctx, cmd.Id, cmd.Timestamp)

	if err != nil {
		return err
	}

	return nil
}

type DeactivateUserGroup struct {
	Id        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

type DeactivateUserGroupHandler decorator.CommandHandler[DeactivateUserGroup]

type deactivateUserGroupHandler struct {
	repo user_groups.Repository
}

func NewDeactivateUserGroupHandler(
	repo user_groups.Repository,
	logger *logger.Logger,
	metricsClient decorator.MetricsClient,
) DeactivateUserGroupHandler {
	if repo == nil {
		panic("nil userGroup repo")
	}

	return decorator.ApplyCommandDecorators[DeactivateUserGroup](
		deactivateUserGroupHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h deactivateUserGroupHandler) Handle(ctx context.Context, cmd DeactivateUserGroup) error {
	err := h.repo.UpdateGroupInActive(ctx, cmd.Id, cmd.Timestamp)

	if err != nil {
		return err
	}

	return nil
}
