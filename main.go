package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"

	"github.com/hinha/coai/config"
	"github.com/hinha/coai/logger"
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
	var g struct {
		errChan <-chan error
	}

	app := fx.New(
		fx.Provide(config.LoadSecret),
		fx.Provide(logger.NewLogger),
		fx.Provide(CreateServer),
		fx.Invoke(Run),
		fx.Extract(&g),
		fx.WithLogger(func(log *logger.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.LogDefault()}
		}),
		fx.StopTimeout(time.Minute),
	)

	startCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <-c
		log.Println("Gracefully shutting down...")
		_ = app.Stop(startCtx)
	}()

	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	select {
	case sg := <-app.Done():
		log.Printf("Received signal: %v", sg)
	case err := <-g.errChan:
		log.Printf("Something server stopped unexpectedly: %v", err)
	}

	os.Exit(0)
}

type Result struct {
	fx.Out

	errChan <-chan error
}

func Run(lc fx.Lifecycle, app *fiber.App, log *logger.Logger, cfg *config.Config) Result {
	defer log.Close()

	errc := make(chan error, 1)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			defer close(errc)
			ctx = context.WithValue(ctx, 1, "TEST")

			log.LogDefault().Info("Running server")
			go StartServer(errc, cfg, app, log)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx = context.WithValue(ctx, 1, "TEST 2")
			log.LogDefault().Info("Stopping server")
			return app.Shutdown()
		},
	})

	return Result{errChan: errc}
}

func CreateServer(cfg *config.Config, log *logger.Logger) *fiber.App {
	log.LogDefault().Info("Starting server")
	app := fiber.New(config.FiberConfig(cfg))
	// middleware
	middlewares.FiberMiddleware(app, cfg, log)

	// routes
	routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	routes.PrivateRoutes(app) // Register a private routes
	routes.NotFoundRoute(app) // Register route for 404 Error.x

	procs := strconv.Itoa(runtime.GOMAXPROCS(0))
	if !app.Config().Prefork {
		procs = "1"
	}

	log.LogDefault().Info("Preparing server",
		zap.String("addr", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)),
		zap.Uint32("handlers", app.HandlersCount()),
		zap.String("num_cpu", procs),
		zap.String("network", app.Config().Network),
		zap.String("mode", string(cfg.Server.Mode)),
		zap.Int("processes", os.Getpid()),
	)
	return app
}

// StartServer func for starting a server.
func StartServer(errChan chan error, cfg *config.Config, a *fiber.App, log *logger.Logger) {
	// Run server.
	go func() {
		if err := a.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
			log.LogDefault().Error("Oops... Server is not running!", zap.Error(err))
			errChan <- err
		}
	}()
}
