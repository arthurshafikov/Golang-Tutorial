package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Queue     amqp.Queue
	QueueName string
	QueueCh   *amqp.Channel
}

func New(ctx context.Context, url string, queueName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel %w", err)
	}
	go func() {
		<-ctx.Done()

		conn.Close()
		ch.Close()
	}()

	queue, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue %w", err)
	}

	return &RabbitMQ{
		Queue:     queue,
		QueueName: queueName,
		QueueCh:   ch,
	}, nil
}

func (r *RabbitMQ) SendToQueue(message []byte) error {
	err := r.QueueCh.Publish(
		"",          // exchange
		r.QueueName, // routing key q.Name
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)

	return err
}

func (r *RabbitMQ) GetDeliveryChannel() (<-chan amqp.Delivery, error) {
	deliveryCh, err := r.QueueCh.Consume(
		r.QueueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	return deliveryCh, err
}
