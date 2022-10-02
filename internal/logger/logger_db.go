package logger

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/utils"
	"time"

	zap_logger "github.com/hinha/zap-logger"
	gorm_logger "gorm.io/gorm/logger"
)

type gormLogger struct {
	logger *zap_logger.ZapLogger

	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  gorm_logger.LogLevel
}

func newGorm(logger *zap_logger.Logger, logLevel int) *gormLogger {
	return &gormLogger{
		logger:        logger,
		LogLevel:      gorm_logger.LogLevel(logLevel),
		SlowThreshold: 200 * time.Millisecond,
	}
}

func (log *gormLogger) LogMode(level gorm_logger.LogLevel) gorm_logger.Interface {
	copy := *log
	copy.LogLevel = level
	return &copy
}

func (log *gormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if log.LogLevel >= gorm_logger.Info {
		log.logger.InfoCtx(ctx, s, zap.Any("tracer", append([]interface{}{utils.FileWithLineNum()}, i...)))
	}
}

func (log *gormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if log.LogLevel >= gorm_logger.Warn {
		log.logger.WarnCtx(ctx, s, zap.Any("tracer", append([]interface{}{utils.FileWithLineNum()}, i...)))
	}
}

func (log *gormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if log.LogLevel >= gorm_logger.Error {
		log.logger.ErrorCtx(ctx, s, zap.Any("tracer", append([]interface{}{utils.FileWithLineNum()}, i...)))
	}
}

func (log *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if log.LogLevel <= gorm_logger.Silent {
		return
	}
	elapsed := time.Since(begin)

	var fields []zap.Field
	fields = append(fields,
		zap.String("src", utils.FileWithLineNum()),
		zap.String("time", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)))

	sql, rows := fc()
	if rows == -1 {
		fields = append(fields, zap.String("row", "-"), zap.String("sql", sql))
	} else {
		fields = append(fields, zap.Int64("row", rows), zap.String("sql", sql))
	}

	switch {
	case err != nil && log.LogLevel >= gorm_logger.Error && (!errors.Is(err, gorm_logger.ErrRecordNotFound) || !log.IgnoreRecordNotFoundError):
		log.logger.ErrorCtx(ctx, "trace", append(fields, zap.Error(err))...)
	case elapsed > log.SlowThreshold && log.SlowThreshold != 0 && log.LogLevel >= gorm_logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", log.SlowThreshold)
		log.logger.WarnCtx(ctx, "trace", append(fields, zap.String("slow", slowLog))...)
	case log.LogLevel == gorm_logger.Info:
		fmt.Println(sql)
		log.logger.InfoCtx(ctx, "trace")
	}
}
