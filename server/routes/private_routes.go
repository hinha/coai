package routes

import (
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	userHandlers "github.com/hinha/coai/modules/user/handlers"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	route.Use(
		cors.New(),
		otelfiber.Middleware(a.Config().AppName),
	)

	// handler Auth
	userHandler := userHandlers.NewAuthHandler()
	userGroup := route.Group("/auth")
	userGroup.Get("/login", userHandler.SignIn).Name("authLogin")

	route.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{})
	}).Name("Book")
}
