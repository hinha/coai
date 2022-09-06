package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"

	"github.com/hinha/coai/config"
	"github.com/hinha/coai/middlewares"
	"github.com/hinha/coai/routes"
)

// @title API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email your@mail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.LoadSecret("./config/config.yml")
	if err != nil {
		panic(err)
	}

	app := fiber.New(config.FiberConfig(cfg))

	// middleware
	middlewares.FiberMiddleware(app, cfg)

	// routes
	routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	routes.PrivateRoutes(app) // Register a private routes
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with or without graceful shutdown).
	if cfg.Server.Mode == "dev" {
		StartServer(cfg, app)
	} else {
		StartServerWithGracefulShutdown(cfg, app)
	}
}

// StartServerWithGracefulShutdown function for starting server with a graceful shutdown.
func StartServerWithGracefulShutdown(cfg *config.Config, a *fiber.App) {
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := a.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	// Run server.
	if err := a.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}

// StartServer func for starting a server.
func StartServer(cfg *config.Config, a *fiber.App) {
	// Run server.
	if err := a.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
