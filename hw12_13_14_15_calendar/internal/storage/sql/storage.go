package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //nolint:gci
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

const (
	EventsTable = "events"
)

type Storage struct {
	db     *sqlx.DB
	dsn    string
	Events []storage.Event
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Connect("postgres", s.dsn)
	go func() {
		<-ctx.Done()
		s.Close()
	}()
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	return err
}

func (s *Storage) Add(event storage.Event) error {
	if s.checkIfStartDateTimeIsBusy(event) {
		return storage.ErrDateBusy
	}
	_, err := s.db.Exec(
		fmt.Sprintf("INSERT INTO %s (title, descr, owner, start_at, end_at, send_notification_at) VALUES($1, $2, $3, $4, $5, $6);", EventsTable),
		event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Change(event storage.Event) error {
	if s.checkIfStartDateTimeIsBusy(event) {
		return storage.ErrDateBusy
	}
	res, err := s.db.Exec(
		fmt.Sprintf("UPDATE %s SET (title, descr, owner, start_at, end_at, send_notification_at) = ($1, $2, $3, $4, $5, $6) WHERE id=$7;", EventsTable),
		event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt, event.ID,
	)
	if err != nil {
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Get(event storage.Event) (storage.Event, error) {
	err := s.db.Get(&event, fmt.Sprintf("SELECT * FROM %s WHERE id=$1 LIMIT 1;", EventsTable), event.ID)
	if err != nil {
		return storage.Event{}, storage.ErrNotFound
	}

	return event, nil
}

func (s *Storage) Delete(event storage.Event) error {
	res, err := s.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id=$1;", EventsTable), event.ID)
	if err != nil {
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListEventsOnADay(t time.Time) (storage.EventsSlice, error) {
	events := storage.EventsSlice{}

	err := s.db.Select(&events, fmt.Sprintf("SELECT * FROM %s WHERE DATE(start_at)=$1;", EventsTable),
		t.Format(storage.RequestDateFormat))

	return events, err
}

func (s *Storage) ListEventsOnARange(t time.Time, tMax time.Time) (storage.EventsSlice, error) {
	events := storage.EventsSlice{}

	err := s.db.Select(&events, fmt.Sprintf("SELECT * FROM %s WHERE start_at BETWEEN $1 AND $2;", EventsTable),
		t.Format(storage.RequestDateTimeFormat), tMax.Format(storage.RequestDateTimeFormat))

	return events, err
}

func (s *Storage) checkIfStartDateTimeIsBusy(event storage.Event) bool {
	res, err := s.db.Exec(fmt.Sprintf("SELECT * FROM %s WHERE start_at=$1;", EventsTable), event.StartAt)
	if err != nil {
		return true
	}
	if n, err := res.RowsAffected(); err != nil || n != 0 {
		return true
	}

	return false
}
