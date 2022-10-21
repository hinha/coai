package application

import (
	"github.com/hinha/coai/core/users/application/command"
	"github.com/hinha/coai/core/users/application/query"
)

type Application struct {
	Queries  Queries
	Commands Commands
}

type Commands struct {
	Register        command.RegisterUserHandler
	CreateGroup     command.CreateUserGroupHandler
	UpdateGroup     command.UpdateUserGroupHandler
	DeleteGroup     command.DeleteUserGroupHandler
	ActivateGroup   command.ActivateUserGroupHandler
	DeactivateGroup command.DeactivateUserGroupHandler
	PingConnection  command.PingHandler
}

type Queries struct {
	AllUsers          query.AllUserHandler
	GetGroup          query.GetGroupHandler
	GroupAvailability query.GroupAvailabilityHandler
}
