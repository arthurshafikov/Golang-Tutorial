package app

import (
	"log"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/config"
	"github.com/thewolf27/hw12_13_14_15_calendar/pkg/rabbitmq"
	"golang.org/x/net/context"
)

type Sender struct {
	Logger   Logger
	Config   config.Config
	RabbitMQ *rabbitmq.RabbitMQ
}

func NewSender(logg Logger, config config.Config) *Sender {
	return &Sender{
		Logger: logg,
		Config: config,
	}
}

func (s *Sender) Run(ctx context.Context) {
	if err := s.сonnect(ctx); err != nil {
		s.Logger.Error(err.Error())
		return
	}

	deliveryCh, err := s.RabbitMQ.GetDeliveryChannel()
	if err != nil {
		s.Logger.Error(err.Error())
		return
	}

	go func() {
		for delivery := range deliveryCh {
			log.Printf("Received a message: %s", delivery.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-ctx.Done()
}

func (s *Sender) сonnect(ctx context.Context) error {
	rabbitMQ, err := rabbitmq.New(ctx, s.Config.RabbitMq.URL, QueueName)
	if err != nil {
		return err
	}
	s.RabbitMQ = rabbitMQ

	return nil
}
