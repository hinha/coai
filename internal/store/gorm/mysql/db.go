package mysql

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/hinha/coai/config"
	"github.com/hinha/coai/internal/logger"
)

var Module = fx.Module(
	"mysql",
	fx.Provide(func(config *config.Config) Config {
		return Config{
			Host:     config.DB.Drivers.Mysql.Host,
			Port:     config.DB.Drivers.Mysql.Port,
			User:     config.DB.Drivers.Mysql.User,
			Password: config.DB.Drivers.Mysql.Pass,
			DBName:   config.DB.Drivers.Mysql.DbName,
		}
	}),
	fx.Provide(New),
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func (d Config) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", d.User, d.Password, d.Host, d.Port, d.DBName)
}

type DB struct {
	Gorm *gorm.DB
	sql  *sql.DB
}

func New(config Config, logger *logger.Logger) *DB {
	dsn := config.String()
	db, err := gorm.Open(mysql.New(
		mysql.Config{DSN: dsn}),
		&gorm.Config{
			Logger: logger.Gorm(),
		})
	if err != nil {
		logger.Console().Fatal("unable to open db connection", zap.Error(err))
	}

	err = db.Use(otelgorm.NewPlugin(otelgorm.WithDBName(config.DBName)))
	if err != nil {
		logger.Console().Fatal("failed to set gorm plugin for opentelemetry", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Console().Fatal("failed to get sql db", zap.Error(err))
	}

	if err := sqlDB.Ping(); err != nil {
		logger.Console().Fatal("failed to ping sql db", zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(200)
	logger.Console().Debug("finish initialize db")

	return &DB{
		Gorm: db,
		sql:  sqlDB,
	}
}

func (c *DB) Ping() error {
	return c.sql.Ping()
}

func (c *DB) Close() error {
	return c.sql.Close()
}
