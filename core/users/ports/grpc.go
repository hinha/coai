package ports

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	app "github.com/hinha/coai/core/users/application"
	"github.com/hinha/coai/core/users/application/command"
	"github.com/hinha/coai/core/users/application/query"
	"github.com/hinha/coai/internal/genproto/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var (
	loc, _     = time.LoadLocation("Asia/Jakarta")
	timeServer time.Time
)

func init() {
	timeServer = time.Now().In(loc)
}

type GrpcServer struct {
	app app.Application
}

func NewGrpcServer(application app.Application) users.UserGroupServiceServer {
	return &GrpcServer{app: application}
}

func (g *GrpcServer) CreateGroup(ctx context.Context, request *users.CreateGroupRequest) (*emptypb.Empty, error) {
	liveTime, err := protoTimestampToTime(request.Time)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = g.app.Commands.CreateGroup.Handle(ctx, command.CreateUserGroup{
		Name:      request.Name,
		Active:    request.Active,
		Timestamp: liveTime,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g *GrpcServer) UpdateGroup(ctx context.Context, request *users.UpdateGroupRequest) (*emptypb.Empty, error) {
	liveTime, err := protoTimestampToTime(request.Time)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = g.app.Commands.UpdateGroup.Handle(ctx, command.UpdateUserGroup{
		Id:        request.Id,
		Active:    request.Active,
		Timestamp: liveTime,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g *GrpcServer) DeleteGroup(ctx context.Context, request *users.DeleteGroupRequest) (*emptypb.Empty, error) {
	liveTime, err := protoTimestampToTime(request.Time)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = g.app.Commands.DeleteGroup.Handle(ctx, command.DeleteUserGroup{
		Id:        request.Id,
		Timestamp: liveTime,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

//func (g *GrpcServer) GetGroup(ctx context.Context, request *users.GetGroupRequest) (*users.GetGroupResponse, error) {
//
//}

func (g *GrpcServer) ActivateGroup(ctx context.Context, request *users.ActivateGroupRequest) (*emptypb.Empty, error) {
	liveTime, err := protoTimestampToTime(request.Time)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = g.app.Commands.ActivateGroup.Handle(ctx, command.ActivateUserGroup{
		Id:        request.Id,
		Timestamp: liveTime,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g *GrpcServer) GetGroup(ctx context.Context, request *users.GetGroupRequest) (*users.GetGroupResponse, error) {
	group, err := g.app.Queries.GetGroup.Handle(ctx, query.GetGroup{
		Name: request.Name,
		Id:   request.Id,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &users.GetGroupResponse{
		Id:        group.ID,
		Name:      group.Name,
		Active:    group.Active,
		CreatedAt: timestamppb.New(group.CreatedAt),
		UpdatedAt: timestamppb.New(group.UpdatedAt),
		DeletedAt: timestamppb.New(group.DeletedAt.Time),
	}, nil
}

func (g *GrpcServer) DeactivateGroup(ctx context.Context, request *users.DeactivateGroupRequest) (*emptypb.Empty, error) {
	liveTime, err := protoTimestampToTime(request.Time)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = g.app.Commands.DeactivateGroup.Handle(ctx, command.DeactivateUserGroup{
		Id:        request.Id,
		Timestamp: liveTime,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g *GrpcServer) IsGroupAvailable(ctx context.Context, request *users.IsGroupAvailableRequest) (*users.IsGroupAvailableResponse, error) {
	isAvailable, err := g.app.Queries.GroupAvailability.Handle(ctx, query.GroupAvailability{Name: request.Name})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &users.IsGroupAvailableResponse{IsAvailable: isAvailable}, nil
}

func protoTimestampToTime(timestamp *timestamp.Timestamp) (time.Time, error) {
	ts := timestamp.AsTime()
	if ts.IsZero() {
		return timeServer, nil
	}

	if ts.Year() > timeServer.Year() {
		return timeServer, errors.New("time request cannot be greater than time server")
	}
	if ts.Year() < timeServer.Year() {
		return timeServer, nil
	}
	return timestamp.AsTime(), nil
}
