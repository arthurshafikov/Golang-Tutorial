package app

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
	"github.com/thewolf27/hw12_13_14_15_calendar/pkg/rabbitmq"
)

const (
	TimeIntervalToScan = time.Second * 3
	EventsTimeToDelete = time.Hour * 24 * 365 // 1 year
	QueueName          = "events"
)

type Scheduler struct {
	Logger      Logger
	Storage     Storage
	RabbitMqURL string
	RabbitMq    *rabbitmq.RabbitMQ
}

type Notification struct {
	ID    int64
	Title string
	Date  string
	Owner int64
}

func NewScheduler(logger Logger, storage Storage, rabbitMqURL string) *Scheduler {
	return &Scheduler{
		Logger:      logger,
		Storage:     storage,
		RabbitMqURL: rabbitMqURL,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	if err := s.connect(ctx); err != nil {
		s.Logger.Error(err.Error())
		return
	}

	intervalTicker := time.NewTicker(TimeIntervalToScan)
	defer intervalTicker.Stop()

OUTER:
	for range intervalTicker.C {
		select {
		case <-ctx.Done():
			break OUTER
		default:
		}

		events, err := s.getEventsThatNeedToBeSend()
		if err != nil {
			s.Logger.Error(err.Error())
			break OUTER
		}

		for _, e := range events {
			notification := Notification{
				ID:    e.ID,
				Title: e.Title,
				Owner: e.Owner,
				Date:  e.StartAt.Format(storage.RequestDateTimeFormat),
			}
			if err := s.setEventIsSentFieldTrue(e); err != nil {
				s.Logger.Error(err.Error())
				return
			}
			if err := s.sendNotificationToQueue(notification); err != nil {
				s.Logger.Error(err.Error())
				return
			}
		}
		s.deleteOldEvents()
	}
}

func (s *Scheduler) connect(ctx context.Context) error {
	rabbitMQ, err := rabbitmq.New(ctx, s.RabbitMqURL, QueueName)
	if err != nil {
		return err
	}
	s.RabbitMq = rabbitMQ

	return nil
}

func (s *Scheduler) getEventsThatNeedToBeSend() (storage.EventsSlice, error) {
	return s.Storage.GetEventsThatNeedToBeSend(time.Now())
}

func (s *Scheduler) deleteOldEvents() {
	eventsToDelete, err := s.Storage.GetEventsWhereEndAtBeforeGivenTimestamp(time.Now().Add(EventsTimeToDelete * -1))
	if err != nil {
		s.Logger.Error(err.Error())
		return
	}

	for _, e := range eventsToDelete {
		err := s.Storage.Delete(e)
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}
	}
}

func (s *Scheduler) sendNotificationToQueue(notification Notification) error {
	message, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	if err := s.RabbitMq.SendToQueue(message); err != nil {
		return err
	}
	log.Printf(" [x] Sent %s", message)

	return nil
}

func (s *Scheduler) setEventIsSentFieldTrue(event storage.Event) error {
	event.IsSent = true
	if _, err := s.Storage.Change(event); err != nil {
		return err
	}

	return nil
}
