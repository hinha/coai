package logger

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"

	"github.com/hinha/coai/config"
	zap_logger "github.com/hinha/zap-logger"
)

var Module = fx.Module(
	"logger",
	fx.Provide(func(config *config.Config) Config {
		return Config{
			Encoding:    config.Log.Output,
			Mode:        string(config.Server.Mode),
			HttpLogPath: config.Log.File.Http,
			GrpcLogPath: config.Log.File.Grpc,
			TimeFormat:  config.Log.TimeFormat,
		}
	}),
	fx.Provide(New),
)

type Config struct {
	Encoding    string
	Mode        string
	HttpLogPath string
	GrpcLogPath string
	TimeFormat  string
}

type Logger struct {
	logHandler *zap_logger.ZapLogger
	logGrpc    *zap_logger.ZapLogger
	logConsole *zap_logger.ZapLogger

	closers []func() error
}

func zapConfig(config Config) zap_logger.Config {
	var zapCfg zap_logger.Config
	if config.Mode == "dev" {
		zapCfg = zap_logger.NewDevelopmentConfig()
	} else {
		zapCfg = zap_logger.NewProductionConfig()
	}
	zapCfg.Encoding = config.Encoding
	zapCfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(config.TimeFormat)
	return zapCfg
}

func New(config Config) *Logger {
	zapConsoleCfg := zapConfig(config)
	zapConsoleCfg.Encoding = "console"
	consoleEncoder := zapcore.NewConsoleEncoder(zapConsoleCfg.EncoderConfig)
	zapConsole := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapConsoleCfg.Level)
	logConsole := zap_logger.New(zapcore.NewTee(zapConsole), zapConsoleCfg)

	zapCfg := zapConfig(config)
	// http config
	zapCfg.Filename = config.HttpLogPath
	logHandler := zap_logger.NewLogger(zapCfg)

	// server config
	zapCfg.Filename = config.GrpcLogPath
	grpcLog := zap_logger.NewLogger(zapCfg)
	return &Logger{
		logHandler: logHandler,
		logGrpc:    grpcLog,
		logConsole: logConsole,
		closers:    []func() error{logHandler.Sync, logConsole.Sync, grpcLog.Sync},
	}
}

func (log *Logger) Console() *zap_logger.ZapLogger {
	return log.logConsole.WithOptions(zap_logger.AddCaller()).Named("console")
}

func (log *Logger) Handler() *zap_logger.ZapLogger {
	return log.logHandler.WithOptions(zap_logger.AddCaller()).Named("handler")
}

func (log *Logger) Core() *zap_logger.ZapLogger {
	return log.logConsole.WithOptions(zap_logger.AddCaller()).Named("core")
}

func (log *Logger) Gorm() *gormLogger {
	logs := log.logConsole.WithOptions(zap_logger.AddCaller()).Named("db")
	if logs.Level() == zapcore.DebugLevel {
		return newGorm(logs, 3) // Warn Level
	}
	return newGorm(logs, 1) // Silent Level
}

func (log *Logger) Grpc() *zap_logger.ZapLogger {
	return log.logGrpc
}

func (log *Logger) Close() error {
	for _, f := range log.closers {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
