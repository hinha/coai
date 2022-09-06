package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/hinha/coai/config"
	"github.com/hinha/coai/logger"
)

func FiberMiddleware(a *fiber.App, cfg *config.Config) {
	fLog, err := os.OpenFile(cfg.Log.File.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	zl := logger.NewZap(fLog, cfg)
	defer func() {
		_ = zl.Sync()
	}()

	a.Use(
		cors.New(),
		NewLogger(Config{AppConfig: cfg, Logger: zl}),
	)
}
