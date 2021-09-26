package storage

import (
	"fmt"
	"time"
)

const (
	DateFormat            = "2006-01-02"
	RequestDateFormat     = "2006-01-02"
	RequestDateTimeFormat = "2006-01-02 15:04:05"
	RequestSuccessMessage = "Success"
)

type EventsSlice []Event

type ListEventsFunction func(time.Time) (EventsSlice, error)

type Event struct {
	ID                 int64
	Title              string
	Descr              string
	Owner              int64
	StartAt            time.Time `db:"start_at"`
	EndAt              time.Time `db:"end_at"`
	SendNotificationAt time.Time `db:"send_notification_at"`
}

var (
	ErrDateBusy = fmt.Errorf("the date and time are busy")
	ErrNotFound = fmt.Errorf("not found")
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

func ParseEvent(input interface{}, strict bool) (Event, error) {
	e := input.(EventRequest)
	var startAt, endAt, sendNotificationAt time.Time
	var err error
	if e.StartAt != "" || strict {
		startAt, err = time.Parse(RequestDateTimeFormat, e.StartAt)
		if err != nil {
			return Event{}, fmt.Errorf("cant parse startAt %w", err)
		}
	}
	if e.EndAt != "" || strict {
		endAt, err = time.Parse(RequestDateTimeFormat, e.EndAt)
		if err != nil {
			return Event{}, fmt.Errorf("cant parse endAt %w", err)
		}
	}
	if e.SendNotificationAt != "" || strict {
		sendNotificationAt, err = time.Parse(RequestDateTimeFormat, e.SendNotificationAt)
		if err != nil {
			return Event{}, fmt.Errorf("cant parse sendNotificationAt %w", err)
		}
	}

	event := Event{
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
