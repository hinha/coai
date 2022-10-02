package server

import (
	"context"
	"fmt"
	"github.com/hinha/coai/internal/store/gorm/mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/hinha/coai/config"
	"github.com/hinha/coai/internal/logger"
	"github.com/hinha/coai/internal/telemetry/exporter"
	"github.com/hinha/coai/internal/telemetry/metric"
	"github.com/hinha/coai/internal/telemetry/trace"
	"github.com/hinha/coai/server/middlewares"
)

var ModuleApp = fx.Module("application",
	fx.Provide(func(cfg *config.Config, logger *logger.Logger) *fiber.App {
		app := fiber.New(config.FiberConfig(cfg))
		app.Use(middlewares.NewLogger(middlewares.Config{
			AppConfig: cfg,
			Logger:    logger.Handler(),
		}))
		return app
	}),
	fx.Provide(NewRouter), // Handle Router
	fx.Provide(NewServer), // Handle server
	fx.Invoke(func(app *fiber.App, router *Router, srv *Server, db *mysql.DB) {
		// close
		srv.closers = append(srv.closers, srv.logger.Close, db.Close)

		// initiate application
		srv.app = app

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

		if srv.config.Otel.Enabled {
			metricExporter := exporter.NewMetricOTLP(srv.config.Otel.Server, srv.logger)
			pusher, pusherCloseFn, err := metric.NewMeterProviderBuilder(srv.logger).
				SetExporter(metricExporter).
				SetHistogramBoundaries([]float64{5, 10, 25, 50, 100, 200, 400, 800, 1000}).
				Build()
			if err != nil {
				srv.logger.Console().Fatal("failed initializing the metric provider", zap.Error(err))
			}
			srv.exporterClosers = append(srv.exporterClosers, pusherCloseFn)
			global.SetMeterProvider(pusher)

			spanExporter := exporter.NewTraceOTLP(srv.config.Otel.Server, srv.logger)
			tracerProvider, tracerProviderCloseFn, err := trace.NewTraceProviderBuilder(srv.config.Server.Name, srv.logger).
				SetExporter(spanExporter).
				Build()
			if err != nil {
				srv.logger.Console().Fatal("failed initializing the tracer provider", zap.Error(err))
			}
			srv.exporterClosers = append(srv.exporterClosers, tracerProviderCloseFn)

			// set global propagator to trace context (the default is no-op).
			otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
			otel.SetTracerProvider(tracerProvider)
		}
	}),
)

var Module = fx.Module("server",
	ModuleRouterSwagger,
	ModuleRouterPublic,
	ModuleRouterPrivate,
	ModuleRouterError,
	ModuleApp,
)

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

	exporterClosers []func(ctx context.Context) error
	closers         []func() error
}

func closer(fn func() error) {
	defer func() {
		_ = fn()
	}()
}

func (s *Server) Close() error {
	for _, f := range s.closers {
		closer(f)
	}

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		cancel()
	}()

	for _, f := range s.exporterClosers {
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
