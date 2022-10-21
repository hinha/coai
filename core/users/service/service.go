package service

import (
	"github.com/hinha/coai/core/users/adapters"
	app "github.com/hinha/coai/core/users/application"
	"github.com/hinha/coai/core/users/application/command"
	"github.com/hinha/coai/core/users/application/query"
	"github.com/hinha/coai/internal/logger"
	"github.com/hinha/coai/internal/metrics"
	"github.com/hinha/coai/internal/store/gorm/mysql"
	"go.uber.org/fx"
)

var Module = fx.Module("application", fx.Provide(NewApplication))

func NewApplication(logger *logger.Logger, db *mysql.DB) app.Application {
	return newApplication(logger, db)
}

func newApplication(logger *logger.Logger, db *mysql.DB) app.Application {
	userRepository := adapters.NewUserMysqlRepository(db)
	userGroupRepository := adapters.NewUserGroupMysqlRepository(db)
	pingRepository := adapters.NewPingMysqlRepository(db)

	metricsClient := metrics.NoOp{}

	return app.Application{
		Commands: app.Commands{
			Register:        command.NewRegisterUserHandler(userRepository, logger, metricsClient),
			PingConnection:  command.NewPingHandler(pingRepository, logger, metricsClient),
			CreateGroup:     command.NewCreateUserGroupHandler(userGroupRepository, logger, metricsClient),
			UpdateGroup:     command.NewUpdateUserGroupHandler(userGroupRepository, logger, metricsClient),
			DeleteGroup:     command.NewDeleteUserGroupHandler(userGroupRepository, logger, metricsClient),
			ActivateGroup:   command.NewActivateUserGroupHandler(userGroupRepository, logger, metricsClient),
			DeactivateGroup: command.NewDeactivateUserGroupHandler(userGroupRepository, logger, metricsClient),
		},
		Queries: app.Queries{
			//AllUsers: query.NewAllUserHandler(repository, logger, metricsClient),
			GetGroup:          query.NewGetGroupHandler(userGroupRepository, logger, metricsClient),
			GroupAvailability: query.NewGroupAvailabilityHandler(userGroupRepository, logger, metricsClient),
		},
	}
}
