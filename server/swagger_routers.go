package server

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"go.uber.org/fx"
)

var ModuleRouterSwagger = fx.Module("router.swagger",
	fx.Invoke(func(a *Router) {
		a.swagger.Get("*", swagger.HandlerDefault)
	}),
)
