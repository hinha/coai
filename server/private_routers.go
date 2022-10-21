package server

import (
	"github.com/hinha/coai/core/users/adapters"
	"go.uber.org/fx"

	"github.com/hinha/coai/core/users/service"
)

// ModuleRouterPrivate Private Handler
var ModuleRouterPrivate = fx.Module("router.private",
	fx.Provide(service.NewApplication),
	fx.Provide(adapters.NewUserHTTP),
	fx.Invoke(func(a *Router, handler *adapters.UserHTTP) {
		group := a.private.Group("/users")
		group.Get("/all", handler.UserAll).Name("userAll")
		group.Get("/register", handler.Register).Name("register")
	}),
)
