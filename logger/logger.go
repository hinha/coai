package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/hinha/coai/config"
)

func NewZap(f *os.File, cfg *config.Config) (zl *zap.Logger) {
	zapConfig := ZapConfig(cfg)
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(cfg.Log.TimeFormat)

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

	return zap.New(core)
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
