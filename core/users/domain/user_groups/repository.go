package user_groups

import (
	"context"
	"github.com/hinha/coai/core/users/domain"
	"time"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package user_groups github.com/hinha/coai
type Repository interface {
	AddUserGroup(ctx context.Context, group *domain.UserGroup, fnInsert func(*domain.UserGroup) (*domain.UserGroup, error)) error
	GetUserGroup(ctx context.Context, id int64, name string) (*domain.UserGroup, error)
	UpdateUserGroup(ctx context.Context, group *domain.UserGroup) error
	DeleteUserGroup(ctx context.Context, group *domain.UserGroup) error
	UpdateGroupActive(ctx context.Context, id int64, time time.Time) error
	UpdateGroupInActive(ctx context.Context, id int64, time time.Time) error
	IsActivateGroup(ctx context.Context, name string) (*domain.UserGroup, error)
}
