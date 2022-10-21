package adapters

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hinha/coai/internal/genproto/users"
)

type UsersGrpc struct {
	client users.UserGroupServiceServer
}

func NewUsersGrpc(client users.UserGroupServiceServer) UsersGrpc {
	return UsersGrpc{client: client}
}

func (s UsersGrpc) CreateGroup(ctx context.Context, name string, active bool, time time.Time) error {
	_, err := s.client.CreateGroup(ctx, &users.CreateGroupRequest{
		Name:   name,
		Active: active,
		Time:   timestamppb.New(time),
	})

	return err
}

func (s UsersGrpc) UpdateGroup(ctx context.Context, id int64, active bool, time time.Time) error {
	_, err := s.client.UpdateGroup(ctx, &users.UpdateGroupRequest{
		Id:     id,
		Active: active,
		Time:   timestamppb.New(time),
	})

	return err
}

func (s UsersGrpc) DeleteGroup(ctx context.Context, id int64, time time.Time) error {
	_, err := s.client.DeleteGroup(ctx, &users.DeleteGroupRequest{
		Id:   id,
		Time: timestamppb.New(time),
	})

	return err
}
