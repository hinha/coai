package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"

	"github.com/hinha/coai/config"
	"github.com/hinha/coai/core/users/ports"
	"github.com/hinha/coai/internal/genproto/users"
	"github.com/hinha/coai/internal/logger"
	"github.com/hinha/coai/internal/store/gorm/mysql"
	"github.com/hinha/coai/internal/telemetry"
	zap_logger "github.com/hinha/zap-logger"
)

var ModuleGrpcServer = fx.Module("grpcServer",
	fx.Provide(func(logger *logger.Logger) *grpc.Server {
		newZap := zap.New(logger.Grpc().Core())
		grpc_zap.ReplaceGrpcLoggerV2(newZap.Named("server"))
		grpcServer := grpc.NewServer(
			grpc_middleware.WithUnaryServerChain(
				grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				grpc_zap.UnaryServerInterceptor(newZap),
				otelgrpc.UnaryServerInterceptor(),
			),
			grpc_middleware.WithStreamServerChain(
				grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				grpc_zap.StreamServerInterceptor(newZap),
				otelgrpc.StreamServerInterceptor(),
			),
		)

		return grpcServer
	}),
	fx.Provide(func(server *grpc.Server) grpc.ServiceRegistrar {
		return server // anonymous interface
	}),
	fx.Provide(ports.NewGrpcServer, ports.NewUserGrpcHealthServer),
	fx.Invoke(users.RegisterUserGroupServiceServer, users.RegisterHealthServer),
	fx.Invoke(RunApp),
)

type registerServer func(server *grpc.Server)

func RunApp(
	lc fx.Lifecycle,
	config *config.Config,
	logger *logger.Logger,
	telemetry *telemetry.OpenTelemetry,
	db *mysql.DB,
	server *grpc.Server) {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			defer logger.Close()
			defer telemetry.Close()
			defer db.Close()

			addr := fmt.Sprintf("%s:%d", config.Server.Grpc.Host, config.Server.Grpc.Port)

			listen, err := net.Listen("tcp", addr)
			if err != nil {
				logger.Console().Fatal("err", zap.Error(err))
			}

			logger.Console().With(zap.String("grpcEndpoint", addr)).Info("Starting: gRPC Listener")

			return server.Serve(listen)
		},
		OnStop: func(ctx context.Context) error {
			server.Stop()
			logger.Console().Info("Server close")
			return nil
		},
	})
}

func RunGRPCServer(config *config.Config, logger *zap_logger.Logger, server registerServer) {
	addr := fmt.Sprintf("%s:%d", config.Server.Grpc.Host, config.Server.Grpc.Port)
	RunGRPCServerOnAddr(addr, logger, server)
}

func RunGRPCServerOnAddr(addr string, logger *zap_logger.Logger, server registerServer) {
	newZap := zap.New(logger.Core())
	grpc_zap.ReplaceGrpcLoggerV2(newZap)

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(newZap),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(newZap),
		),
	)
	server(grpcServer)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		newZap.Fatal("err", zap.Error(err))
	}
	newZap.With(zap.String("grpcEndpoint", addr)).Info("Starting: gRPC Listener")
	newZap.Fatal("listen", zap.Error(grpcServer.Serve(listen)))
}
