package server

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// ModuleRouterError func for describe application Error.
var ModuleRouterError = fx.Module("router.error",
	fx.Invoke(func(a *fiber.App) {
		a.Use(func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "sorry, endpoint is not found",
			})
		})
	}))
