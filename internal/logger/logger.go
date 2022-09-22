package logger

import (
	"os"

	zap_logger "github.com/hinha/zap-logger"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Encoding   string
	Mode       string
	LogPath    string
	TimeFormat string
}

type Logger struct {
	logHandler *zap_logger.ZapLogger
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
	zapCfg.Filename = config.LogPath
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
	logHandler := zap_logger.NewLogger(zapCfg)
	return &Logger{
		logHandler: logHandler,
		logConsole: logConsole,
		closers:    []func() error{logHandler.Sync, logConsole.Sync},
	}
}

func (log *Logger) Console() *zap_logger.ZapLogger {
	return log.logConsole.WithOptions(zap_logger.AddCaller()).Named("console")
}

func (log *Logger) Handler() *zap_logger.ZapLogger {
	return log.logHandler.WithOptions(zap_logger.AddCaller()).Named("handler")
}

func (log *Logger) Close() error {
	for _, f := range log.closers {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
