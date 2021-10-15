package storage

import (
	"time"
)

const (
	RequestDateFormat     = "2006-01-02"
	RequestDateTimeFormat = "2006-01-02 15:04:05"
	RequestSuccessMessage = "Success"
)

type EventsSlice []Event

type Event struct {
	ID                 int64
	Title              string
	Descr              string
	Owner              int64
	StartAt            time.Time `db:"start_at"`
	EndAt              time.Time `db:"end_at"`
	SendNotificationAt time.Time `db:"send_notification_at"`
	IsSent             bool      `db:"is_sent"`
}

type ListEventsFunction func(time.Time) (EventsSlice, error)
