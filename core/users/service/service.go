package service

import (
	"github.com/hinha/coai/core/users/adapters"
	app "github.com/hinha/coai/core/users/application"
	"github.com/hinha/coai/core/users/application/command"
	"github.com/hinha/coai/core/users/application/query"
	"github.com/hinha/coai/internal/logger"
	"github.com/hinha/coai/internal/metrics"
	"github.com/hinha/coai/internal/store/gorm/mysql"
)

func NewApplication(logger *logger.Logger, db *mysql.DB) app.Application {
	return newApplication(logger, db)
}

func newApplication(logger *logger.Logger, db *mysql.DB) app.Application {
	repository := adapters.NewUserMysqlRepository(db)

	metricsClient := metrics.NoOp{}

	return app.Application{
		Commands: app.Commands{
			Register: command.NewRegisterUserHandler(repository, logger, metricsClient),
		},
		Queries: app.Queries{
			AllUsers: query.NewAllUserHandler(repository, logger, metricsClient),
		},
	}
}
