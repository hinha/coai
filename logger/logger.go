package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/hinha/coai/config"
)

type Logger struct {
	*config.Config
	zap     *zap.Logger
	closers []func() error
}

func NewLogger(cfg *config.Config) *Logger {
	f, err := os.OpenFile(cfg.Log.File.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	zapConfig := ZapConfig(cfg)
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(cfg.Log.TimeFormat)
	if cfg.Log.Color {
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	fileEncoder := zapcore.NewJSONEncoder(zapConfig.EncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(zapConfig.EncoderConfig)

	zapConsole := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapConfig.Level)
	zapFile := zapcore.NewCore(fileEncoder, zapcore.AddSync(f), zapConfig.Level)

	var zapCores []zapcore.Core
	switch cfg.Log.Output {
	case config.LogStdout:
		zapCores = append(zapCores, zapConsole, zapFile)
	case config.LogFile:
		zapCores = append(zapCores, zapFile)
	default:
		zapCores = append(zapCores, zapConsole)
	}

	core := zapcore.NewTee(zapCores...)
	zp := zap.New(core)

	return &Logger{Config: cfg, zap: zp, closers: []func() error{f.Close, zp.Sync}}
}

func (l *Logger) LogDefault() *zap.Logger {
	zapConfig := ZapConfig(l.Config)
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(l.Log.TimeFormat)
	if l.Log.Color {
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	consoleEncoder := zapcore.NewConsoleEncoder(zapConfig.EncoderConfig)
	zapConsole := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapConfig.Level)

	return zap.New(zapcore.NewTee(zapConsole))
}

func (l *Logger) Logger() *zap.Logger {
	return l.zap
}

func (l *Logger) Close() error {
	for _, f := range l.closers {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

func ZapConfig(cfg *config.Config) (zc zap.Config) {
	if cfg.Server.Mode == config.Development {
		zc = zap.NewDevelopmentConfig()
	} else {
		zc = zap.NewProductionConfig()
	}

	var encoding string
	if cfg.Log.Format == "text" {
		encoding = "console"
	} else {
		encoding = "json"
	}
	zc.Encoding = encoding
	return
}
