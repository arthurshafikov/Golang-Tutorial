package app

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/thewolf27/hw12_13_14_15_calendar/pkg/rabbitmq"
	"golang.org/x/net/context"
)

type ConsumerMessagesCh chan string

type Sender struct {
	Logger             Logger
	RabbitMqURL        string
	RabbitMQ           *rabbitmq.RabbitMQ
	ConsumerMessagesCh ConsumerMessagesCh
}

func NewSender(logg Logger, rabbitMqURL string) *Sender {
	return &Sender{
		Logger:             logg,
		RabbitMqURL:        rabbitMqURL,
		ConsumerMessagesCh: make(ConsumerMessagesCh),
	}
}

func (s *Sender) Run(ctx context.Context) {
	defer close(s.ConsumerMessagesCh)
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
			s.sendMessageToConsumer(delivery)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-ctx.Done()
}

func (s *Sender) сonnect(ctx context.Context) error {
	rabbitMQ, err := rabbitmq.New(ctx, s.RabbitMqURL, QueueName)
	if err != nil {
		return err
	}
	s.RabbitMQ = rabbitMQ

	return nil
}

func (s *Sender) sendMessageToConsumer(delivery amqp091.Delivery) {
	s.ConsumerMessagesCh <- string(delivery.Body)
}
