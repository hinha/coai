package server

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// Example Public Handler
var ModuleRouterPublic = fx.Module("router.public",
	fx.Invoke(func(a *Router) {
		a.public.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(map[string]interface{}{})
		}).Name("Book")
		a.public.Get("/ab", func(c *fiber.Ctx) error {
			return c.JSON(map[string]interface{}{})
		}).Name("Book")
	}),
)
