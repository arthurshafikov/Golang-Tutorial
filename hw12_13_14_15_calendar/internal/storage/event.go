package storage

import "fmt"

type Interface interface {
	Add(Event) error
	Change(Event) error
	Delete(Event) error
	List() error
}

type Event struct {
	ID                 string
	Title              string
	Descr              string
	Owner              int64
	StartDate          string `db:"start_date"`
	StartTime          string `db:"start_time"`
	EndDate            string `db:"end_date"`
	EndTime            string `db:"end_time"`
	SendNotificationAt string `db:"send_notification_at"`
}

var (
	ErrDateBusy = fmt.Errorf("the date and time are busy")
	ErrNotFound = fmt.Errorf("not found")
)
