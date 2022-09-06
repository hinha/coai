package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/hinha/coai/accounts/handlers"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// handler Auth
	authHandler := handlers.NewAuthHandler()
	authGroup := route.Group("/auth")
	authGroup.Get("/login", authHandler.SignIn).Name("authLogin")

	route.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{})
	}).Name("Book")
}
