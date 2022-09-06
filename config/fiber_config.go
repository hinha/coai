package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func FiberConfig(cfg *Config) fiber.Config {
	return fiber.Config{
		AppName:           cfg.Server.Name,
		ReadTimeout:       time.Second * time.Duration(cfg.Server.Timeout),
		ReduceMemoryUsage: true,
		CaseSensitive:     true,
	}
}
