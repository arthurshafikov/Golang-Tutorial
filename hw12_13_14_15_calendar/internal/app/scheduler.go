package app

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/config"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
	"github.com/thewolf27/hw12_13_14_15_calendar/pkg/rabbitmq"
)

const (
	TimeIntervalToScan = time.Second * 3
	EventsTimeToDelete = time.Hour * 24 * 365 // 1 year
	QueueName          = "events"
)

type Scheduler struct {
	Logger   Logger
	Storage  Storage
	Config   config.Config
	RabbitMQ *rabbitmq.RabbitMQ
}

type Notification struct {
	ID    int64
	Title string
	Date  time.Time
	Owner int64
}

func NewScheduler(logger Logger, storage Storage, config config.Config) *Scheduler {
	return &Scheduler{
		Logger:  logger,
		Storage: storage,
		Config:  config,
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

		go s.deleteOldEvents()

		for _, e := range events {
			notification := Notification{
				ID:    e.ID,
				Title: e.Title,
				Owner: e.Owner,
				Date:  e.StartAt,
			}
			if err := s.sendNotificationToQueue(notification); err != nil {
				s.Logger.Error(err.Error())
				return
			}
			if err := s.setEventIsSentFieldTrue(e); err != nil {
				s.Logger.Error(err.Error())
				return
			}
		}
		log.Println("Ticking every 3 seconds")
	}
}

func (s *Scheduler) connect(ctx context.Context) error {
	rabbitMQ, err := rabbitmq.New(ctx, s.Config.RabbitMq.URL, QueueName)
	if err != nil {
		return err
	}
	s.RabbitMQ = rabbitMQ

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
	if err := s.RabbitMQ.SendToQueue(message); err != nil {
		return err
	}
	log.Printf(" [x] Sent %s", message)

	return nil
}

func (s *Scheduler) setEventIsSentFieldTrue(event storage.Event) error {
	event.IsSent = true
	if err := s.Storage.Change(event); err != nil {
		return err
	}

	return nil
}
