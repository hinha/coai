package server

import (
	"fmt"
	"github.com/hinha/coai/config"
	"github.com/hinha/coai/core/users/service"
	"github.com/hinha/coai/internal/logger"
	"github.com/hinha/coai/internal/server"
	"github.com/hinha/coai/internal/store/gorm/mysql"
	"github.com/hinha/coai/internal/telemetry"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log"
)

func MakeServerCmd() *cli.Command {
	cmd := &cli.Command{
		Name:  "server",
		Usage: "running server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c", "cfg"},
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("file config", c.String("config"))

			app := fx.New(
				fx.Provide(config.LoadSecret),
				logger.Module,
				telemetry.Module,
				mysql.Module,
				service.Module,
				server.ModuleGrpcServer,
				fx.WithLogger(func(log *logger.Logger) fxevent.Logger {
					return log
				}),
			)

			if err := app.Start(c.Context); err != nil {
				log.Fatal(err)
			}

			//cli.ShowCommandHelpAndExit(c, "server", 1)
			return app.Stop(c.Context)
		},
	}

	return cmd
}
