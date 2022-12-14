package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hinha/coai/config"
	"github.com/hinha/coai/logger"
	"github.com/hinha/coai/server"
)

// @title API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email your@mail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go func() {
		sigs := <-ch
		log.Printf("system call:%+v", sigs)
		cancel()
	}()

	app := fx.New(
		fx.Provide(config.LoadSecret, logger.NewLogger),
		fx.Provide(server.NewServer),
		fx.Invoke(server.InitFiber),
		fx.WithLogger(func(log *logger.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.LogDefault()}
		}),
	)

	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}

	select {
	case <-ctx.Done():
		ctxShutDown, cancelShutdown := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancelShutdown()
		if err := app.Stop(ctxShutDown); err != nil {
			log.Fatal(err)
		}
	}
	os.Exit(0)
}
