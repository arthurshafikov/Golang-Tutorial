package server

import (
	"fmt"
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

type EventRequest struct {
	ID                 int64
	Title              string
	Descr              string
	Owner              int64
	StartAt            string
	EndAt              string
	SendNotificationAt string
}

func ParseEvent(input interface{}, requireAllFields bool) (storage.Event, error) {
	e := input.(EventRequest)
	var startAt, endAt, sendNotificationAt time.Time
	var err error
	if e.StartAt != "" || requireAllFields {
		startAt, err = time.Parse(storage.RequestDateTimeFormat, e.StartAt)
		if err != nil {
			return storage.Event{}, fmt.Errorf("cant parse startAt %w", err)
		}
	}
	if e.EndAt != "" || requireAllFields {
		endAt, err = time.Parse(storage.RequestDateTimeFormat, e.EndAt)
		if err != nil {
			return storage.Event{}, fmt.Errorf("cant parse endAt %w", err)
		}
	}
	if e.SendNotificationAt != "" || requireAllFields {
		sendNotificationAt, err = time.Parse(storage.RequestDateTimeFormat, e.SendNotificationAt)
		if err != nil {
			return storage.Event{}, fmt.Errorf("cant parse sendNotificationAt %w", err)
		}
	}

	event := storage.Event{
		ID:                 e.ID,
		Title:              e.Title,
		Descr:              e.Descr,
		Owner:              e.Owner,
		StartAt:            startAt,
		EndAt:              endAt,
		SendNotificationAt: sendNotificationAt,
	}

	return event, nil
}
