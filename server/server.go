package server

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/hinha/coai/config"
	"github.com/hinha/coai/internal/logger"
	"github.com/hinha/coai/internal/telemetry/exporter"
	"github.com/hinha/coai/internal/telemetry/metric"
	"github.com/hinha/coai/internal/telemetry/trace"
	"github.com/hinha/coai/server/middlewares"
	"github.com/hinha/coai/server/routes"
)

func InitFiber(srv *Server) {
	srv.app = fiber.New(config.FiberConfig(srv.config))
	// Register logger
	srv.app.Use(middlewares.NewLogger(middlewares.Config{
		AppConfig: srv.config,
		Logger:    srv.logger.Handler(),
	}))

	routes.SwaggerRoute(srv.app)  // Register a route for API Docs (Swagger).
	routes.PrivateRoutes(srv.app) // Register a private routes
	routes.NotFoundRoute(srv.app) // Register route for 404 Error.x

	procs := strconv.Itoa(runtime.GOMAXPROCS(0))
	if !srv.app.Config().Prefork {
		procs = "1"
	}

	srv.logger.Console().Info("Preparing server",
		zap.String("addr", fmt.Sprintf("%s:%d", srv.config.Server.Host, srv.config.Server.Port)),
		zap.Uint32("handlers", srv.app.HandlersCount()),
		zap.String("num_cpu", procs),
		zap.String("network", srv.app.Config().Network),
		zap.String("mode", string(srv.config.Server.Mode)),
		zap.Int("processes", os.Getpid()),
	)

	// define
	if srv.config.Otel.Enabled {
		metricExporter := exporter.NewMetricOTLP(srv.config.Otel.Server, srv.logger)
		pusher, pusherCloseFn, err := metric.NewMeterProviderBuilder(srv.logger).
			SetExporter(metricExporter).
			SetHistogramBoundaries([]float64{5, 10, 25, 50, 100, 200, 400, 800, 1000}).
			Build()
		if err != nil {
			srv.logger.Console().Error("failed initializing the metric provider", zap.Error(err))
		}
		srv.Closers = append(srv.Closers, pusherCloseFn)
		global.SetMeterProvider(pusher)

		spanExporter := exporter.NewTraceOTLP(srv.config.Otel.Server, srv.logger)
		tracerProvider, tracerProviderCloseFn, err := trace.NewTraceProviderBuilder(srv.config.Server.Name, srv.logger).
			SetExporter(spanExporter).
			Build()
		if err != nil {
			srv.logger.Console().Error("failed initializing the tracer provider", zap.Error(err))
		}
		srv.Closers = append(srv.Closers, tracerProviderCloseFn)

		// set global propagator to trace context (the default is no-op).
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		otel.SetTracerProvider(tracerProvider)
	}
}

// NewServer ...
func NewServer(lc fx.Lifecycle, cfg *config.Config, logger *logger.Logger) *Server {
	s := &Server{
		logger: logger,
		config: cfg,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go s.Start()
			logger.Console().Info("Server serving", zap.Int("port", cfg.Server.Port))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Console().Info("Server stopped")
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
			s.logger.Console().Fatal("Oops... Server is not running!", zap.Error(err))
		}
	}()
}
