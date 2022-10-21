package ports

import (
	"context"
	app "github.com/hinha/coai/core/users/application"
	"github.com/hinha/coai/core/users/application/command"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"

	"github.com/hinha/coai/internal/genproto/users"
)

type GrpcHealthServerUser struct {
	app app.Application
}

func NewUserGrpcHealthServer(application app.Application) users.HealthServer {
	return &GrpcHealthServerUser{app: application}
}

func (g *GrpcHealthServerUser) Check(ctx context.Context, request *users.HealthCheckRequest) (*users.HealthCheckResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GrpcHealthServerUser) Watch(_ *emptypb.Empty, server users.Health_WatchServer) error {
	var (
		servingCount    int
		notServingCount int
		downTime        float64
		upTime          float64
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response := new(users.HealthCheckResponse)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	done := make(chan bool)
	defer close(done)
	for {
		select {
		case <-done:
			return nil
		case tk := <-ticker.C:
			if err := g.app.Commands.PingConnection.Handle(ctx, command.Ping{StateUP: servingCount, StateDOWN: notServingCount}); err != nil {
				if servingCount > 0 {
					servingCount = 0
					upTime = 0
				}

				notServingCount++
				downTime += float64(time.Since(tk)) / float64(time.Millisecond)
				response.Status = users.HealthCheckResponse_NOT_SERVING
				response.DowntimeDuration = float32(downTime)
			} else {
				// reset downTime and notServing
				if notServingCount > 0 {
					notServingCount = 0
					downTime = 0
				}

				servingCount++
				upTime += float64(time.Since(tk)) / float64(time.Millisecond)
				response.UptimeDuration = float32(upTime)
				response.Status = users.HealthCheckResponse_SERVING
			}

			if err := server.Send(response); err != nil {
				ticker.Stop()
				done <- true
			}
		}
	}
}
