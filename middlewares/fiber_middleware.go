package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/hinha/coai/config"
	"github.com/hinha/coai/logger"
)

func FiberMiddleware(a *fiber.App, cfg *config.Config, zap *logger.Logger) {
	a.Use(
		cors.New(),
		NewLogger(Config{AppConfig: cfg, Logger: zap.Logger()}),
	)
}
