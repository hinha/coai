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
	Register command.RegisterUserHandler
}

type Queries struct {
	AllUsers query.AllUserHandler
}
