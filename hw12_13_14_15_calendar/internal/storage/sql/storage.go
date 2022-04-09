package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/arthurshafikov/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //nolint:gci
)

const (
	EventsTable                   = "events"
	EventStartAtColumn            = "start_at"
	EventEndAtColumn              = "end_at"
	EventSendNotificationAtColumn = "send_notification_at"
	EventIsSentColumn             = "is_sent"
	EventsColumns                 = "title, descr, owner, start_at, end_at, send_notification_at, is_sent"
	EventsIsSentFalse             = "f"
	EventsIsSentTrue              = "t"
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
	go func() {
		<-ctx.Done()
		s.Close()
	}()
	if err != nil {
		panic(err)
	}
	s.db = db
	return nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	return err
}

func (s *Storage) Add(event storage.Event) (int64, error) {
	event.ID = 0
	if s.checkIfStartAtIsBusy(event) {
		return 0, storage.ErrStartAtBusy
	}

	var lastInsertID int64
	err := s.db.QueryRowx(
		fmt.Sprintf("INSERT INTO %s (%s) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id;", EventsTable, EventsColumns),
		event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt, EventsIsSentFalse,
	).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("cannot INSERT event")
	}

	return lastInsertID, nil
}

func (s *Storage) Change(event storage.Event) (int64, error) {
	if s.checkIfStartAtIsBusy(event) {
		return 0, storage.ErrStartAtBusy
	}

	var lastInsertID int64
	err := s.db.QueryRowx(
		fmt.Sprintf("UPDATE %s SET (%s) = ($1, $2, $3, $4, $5, $6, $7) WHERE id=$8 RETURNING id;", EventsTable, EventsColumns),
		event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt, event.IsSent, event.ID,
	).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("cannot UPDATE event")
	}

	return lastInsertID, nil
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

func (s *Storage) ListEventsOnADay(date time.Time) (storage.EventsSlice, error) {
	events := storage.EventsSlice{}

	err := s.db.Select(&events, fmt.Sprintf("SELECT * FROM %s WHERE DATE(%s)=$1;", EventsTable, EventStartAtColumn),
		date.Format(storage.RequestDateFormat))

	return events, err
}

func (s *Storage) ListEventsOnARange(rangeStartTime, rangeEndTime time.Time) (storage.EventsSlice, error) {
	events := storage.EventsSlice{}

	err := s.db.Select(&events, fmt.Sprintf("SELECT * FROM %s WHERE %s BETWEEN $1 AND $2;", EventsTable, EventStartAtColumn),
		rangeStartTime.Format(storage.RequestDateTimeFormat), rangeEndTime.Format(storage.RequestDateTimeFormat))

	return events, err
}

func (s *Storage) GetEventsThatNeedToBeSend(timeTo time.Time) (storage.EventsSlice, error) {
	events := storage.EventsSlice{}

	err := s.db.Select(&events, fmt.Sprintf("SELECT * FROM %s WHERE %s <= $1 AND %s = 'f';", EventsTable, EventSendNotificationAtColumn, EventIsSentColumn),
		timeTo.Format(storage.RequestDateTimeFormat))

	return events, err
}

func (s *Storage) GetEventsWhereEndAtBeforeGivenTimestamp(timeTo time.Time) (storage.EventsSlice, error) {
	events := storage.EventsSlice{}

	err := s.db.Select(&events, fmt.Sprintf("SELECT * FROM %s WHERE %s <= $1;", EventsTable, EventEndAtColumn),
		timeTo.Format(storage.RequestDateTimeFormat))

	return events, err
}

func (s *Storage) ResetDB() error {
	_, err := s.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", EventsTable))

	return err
}

func (s *Storage) checkIfStartAtIsBusy(event storage.Event) bool {
	res, err := s.db.Exec(fmt.Sprintf("SELECT * FROM %s WHERE %s=$1 AND id!=$2;", EventsTable, EventStartAtColumn), event.StartAt, event.ID)
	if err != nil {
		return true
	}
	if n, err := res.RowsAffected(); err != nil || n != 0 {
		return true
	}

	return false
}
