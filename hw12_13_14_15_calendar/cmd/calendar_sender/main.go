package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/app"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/launch"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config, logg := launch.InitializateConfigAndLoggerFromFlags()

	sender := app.NewSender(logg, config.RabbitMq.URL)

	go func() {
		for mes := range sender.ConsumerMessagesCh {
			log.Printf("Received a message: %s", mes)
		}
	}()

	sender.Run(ctx)
}
