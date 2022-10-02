package server

import (
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Router struct {
	public  fiber.Router
	private fiber.Router
	swagger fiber.Router
}

// NewRouter func for describe multiple group
func NewRouter(a *fiber.App) *Router {
	public := a.Group("/public")
	private := a.Group("/api/v1")
	private.Use(
		cors.New(),
		otelfiber.Middleware(a.Config().AppName),
	)
	swagger := a.Group("/swagger")
	return &Router{public: public, private: private, swagger: swagger}
}
