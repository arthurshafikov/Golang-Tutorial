package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/app"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/launch"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage/resolver"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config, logg := launch.InitializateConfigAndLoggerFromFlags()

	storage, err := resolver.ResolveStorage(ctx, config)
	if err != nil {
		logg.Error(err.Error())
		return
	}

	scheduler := app.NewScheduler(logg, storage.(app.Storage), config.RabbitMq.URL)
	scheduler.Run(ctx)
}
