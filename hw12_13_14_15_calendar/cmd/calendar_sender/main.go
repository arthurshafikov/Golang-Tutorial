package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/app"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/launch"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config, logg := launch.Initializate()

	sender := app.NewSender(logg, config)
	sender.Run(ctx)
}
