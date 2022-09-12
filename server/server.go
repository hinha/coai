package server

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hinha/coai/config"
	"github.com/hinha/coai/logger"
	"github.com/hinha/coai/server/middlewares"
	"github.com/hinha/coai/server/routes"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"runtime"
	"strconv"
	"time"
)

func InitFiber(srv *Server) {
	srv.app = fiber.New(config.FiberConfig(srv.config))
	// Register logger
	srv.app.Use(middlewares.NewLogger(middlewares.Config{
		AppConfig: srv.config,
		Logger:    srv.logger.Logger(),
	}))

	routes.SwaggerRoute(srv.app)  // Register a route for API Docs (Swagger).
	routes.PrivateRoutes(srv.app) // Register a private routes
	routes.NotFoundRoute(srv.app) // Register route for 404 Error.x

	procs := strconv.Itoa(runtime.GOMAXPROCS(0))
	if !srv.app.Config().Prefork {
		procs = "1"
	}

	srv.logger.LogDefault().Info("Preparing server",
		zap.String("addr", fmt.Sprintf("%s:%d", srv.config.Server.Host, srv.config.Server.Port)),
		zap.Uint32("handlers", srv.app.HandlersCount()),
		zap.String("num_cpu", procs),
		zap.String("network", srv.app.Config().Network),
		zap.String("mode", string(srv.config.Server.Mode)),
		zap.Int("processes", os.Getpid()),
	)
}

// NewServer ...
func NewServer(lc fx.Lifecycle, cfg *config.Config, logger *logger.Logger, options ...Option) *Server {
	s := &Server{
		logger: logger,
		config: cfg,
	}
	// TODO: move handling inside fx injector
	options = append(options, WithTelemetry(cfg.Server.Name, "0.0.0.0:4317"))

	// loop through our parsing options and apply them
	for _, option := range options {
		option(s)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go s.Start()
			logger.Logger().Info("Server serving", zap.Int("port", cfg.Server.Port))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Logger().Info("Server stopped")
			return s.Close()
		},
	})
	return s
}

type Server struct {
	app    *fiber.App
	logger *logger.Logger
	config *config.Config

	Closers []func(ctx context.Context) error
}

func (s *Server) Close() error {
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		cancel()
	}()

	for _, f := range s.Closers {
		if err := f(ctxShutDown); err != nil {
			return err
		}
	}

	return s.app.Shutdown()
}

func (s *Server) Start() {
	go func() {
		if err := s.app.Listen(fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)); err != nil {
			s.logger.LogDefault().Fatal("Oops... Server is not running!", zap.Error(err))
		}
	}()
}
