package sqlstorage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //nolint:gci
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db  *sqlx.DB
	dsn string
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Connect("postgres", s.dsn)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	err := s.db.Close()
	return err
}

func (s *Storage) Add(event storage.Event) error {
	if s.checkIfStartDateTimeIsBusy(event) {
		return storage.ErrDateBusy
	}
	res, err := s.db.Exec(
		`INSERT INTO events(title, descr, owner, start_date, end_date, send_notification_at) VALUES($1, $2, $3, $4, $5, $6);`,
		event.Title, event.Descr, event.Owner, event.StartDate, event.EndDate, event.SendNotificationAt,
	)
	if err != nil {
		return err
	}
	if _, err := res.LastInsertId(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Change(event storage.Event) error {
	if s.checkIfStartDateTimeIsBusy(event) {
		return storage.ErrDateBusy
	}
	res, err := s.db.Exec(
		`UPDATE events
		SET (title, descr, owner, start_date, end_date, send_notification_at)
		VALUES($1, $2, $3, $4, $5, $6)
		WHERE id=$7;`,
		event.Title, event.Descr, event.Owner, event.StartDate, event.EndDate, event.SendNotificationAt, event.ID,
	)
	if err != nil {
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Get(id string) (storage.Event, error) {
	event := storage.Event{}
	err := s.db.Get(&event, "SELECT * FROM events WHERE id=$1;", id)
	if err != nil {
		return storage.Event{}, storage.ErrNotFound
	}

	return event, nil
}

func (s *Storage) Delete(event storage.Event) error {
	res, err := s.db.Exec(
		`DELETE FROM events WHERE id=$1;`, event.ID,
	)
	if err != nil {
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) List() error {
	var events []storage.Event
	err := s.db.Select(&events, "SELECT * FROM events;")
	if err != nil {
		return err
	}
	for _, e := range events {
		fmt.Printf("%#v\n", e)
	}
	return nil
}

func (s *Storage) checkIfStartDateTimeIsBusy(event storage.Event) bool {
	res, err := s.db.Exec("SELECT * FROM events WHERE start_date=$1 AND start_time=$2;", event.StartDate, event.StartTime)
	if err != nil {
		return true
	}
	n, err := res.RowsAffected()
	if err != nil || n != 0 {
		return true
	}
	return false
}
