package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/app"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/launch"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/server/grpc/api"
	internalhttp "github.com/thewolf27/hw12_13_14_15_calendar/internal/server/http"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage/resolver"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config, logg := launch.Initializate()

	storage, err := resolver.ResolveStorage(ctx, config)
	if err != nil {
		logg.Error(err.Error())
		return
	}

	calendar := app.New(logg, storage.(app.Storage))

	go api.RunGrpcServer(ctx, config.GrpcServer.Host, config.GrpcServer.Port, calendar, logg)

	server := internalhttp.NewServer(logg, calendar, config.HTTPServer.Host, config.HTTPServer.Port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		return
	}
}
