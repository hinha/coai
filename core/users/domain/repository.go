package domain

import "context"

type Repository interface {
	AddUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, userID int64) (*User, error)
}
