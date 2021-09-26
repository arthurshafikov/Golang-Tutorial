package sqlstorage

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

var events = storage.EventsSlice{
	storage.Event{
		ID:      1,
		Title:   "test",
		Descr:   "testdes2",
		Owner:   42,
		StartAt: time.Date(2019, 04, 31, 05, 13, 58, 10, time.Time{}.UTC().Location()),
		EndAt:   time.Date(2019, 05, 06, 05, 13, 58, 10, time.Time{}.UTC().Location()),
	},
	storage.Event{
		ID:      2,
		Title:   "test2",
		Descr:   "testdes2",
		Owner:   13,
		StartAt: time.Date(2020, 04, 31, 05, 13, 58, 10, time.Time{}.UTC().Location()),
		EndAt:   time.Date(2020, 05, 06, 05, 13, 58, 10, time.Time{}.UTC().Location()),
	},
}

var eventColumns = []string{
	"id", "title", "descr", "owner", "start_at", "end_at", "send_notification_at",
}

func newSQLStorageMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, Storage) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mockStorage := Storage{
		dsn: "test dsn",
		db:  sqlxDB,
	}

	return mockDB, mock, mockStorage
}

func TestAddEvent(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	event := events[0]

	mock.ExpectExec(fmt.Sprintf("SELECT (.+) FROM %s", EventsTable)).WithArgs(
		event.StartAt,
	).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", EventsTable)).WithArgs(
		event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err := mockStorage.Add(event)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestChangeEvent(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	event := events[0]

	mock.ExpectExec(fmt.Sprintf("SELECT (.+) FROM %s", EventsTable)).WithArgs(
		event.StartAt,
	).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(fmt.Sprintf("UPDATE %s SET", EventsTable)).WithArgs(
		event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt, event.ID,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err := mockStorage.Change(event)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetEvent(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	event := events[0]

	mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s", EventsTable)).WithArgs(event.ID).WillReturnRows(
		sqlmock.NewRows(eventColumns).AddRow(
			event.ID, event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt,
		),
	)

	resEvent, err := mockStorage.Get(event)
	require.NoError(t, err)
	require.Equal(t, event, resEvent)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestDeleteEvent(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	event := events[0]

	mock.ExpectExec(fmt.Sprintf("DELETE FROM %s", EventsTable)).WithArgs(
		event.ID,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err := mockStorage.Delete(event)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestListTodayEvents(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	event1 := events[0]
	event2 := events[0]

	mock.ExpectQuery(
		fmt.Sprintf("SELECT (.+) FROM %s WHERE DATE\\(start_at\\)", EventsTable)).
		WithArgs(event1.StartAt.Format(storage.RequestDateFormat)).
		WillReturnRows(
			sqlmock.NewRows(eventColumns).AddRow(
				event1.ID, event1.Title, event1.Descr, event1.Owner, event1.StartAt, event1.EndAt, event1.SendNotificationAt,
			).AddRow(
				event2.ID, event2.Title, event2.Descr, event2.Owner, event2.StartAt, event2.EndAt, event2.SendNotificationAt,
			),
		)

	result, err := mockStorage.ListEventsOnADay(event1.StartAt)
	require.NoError(t, err)
	expected := storage.EventsSlice{event1, event2}
	require.Equal(t, expected, result)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestListWeekEvents(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	event1 := events[0]
	event2 := events[0]
	event2.StartAt = event2.StartAt.Add(time.Hour * 24 * 2)
	tMax := event1.StartAt.Add(7 * time.Hour * 24)

	mock.ExpectQuery(
		fmt.Sprintf("SELECT (.+) FROM %s WHERE start_at BETWEEN", EventsTable)).
		WithArgs(event1.StartAt.Format(storage.RequestDateTimeFormat), tMax.Format(storage.RequestDateTimeFormat)).
		WillReturnRows(
			sqlmock.NewRows(eventColumns).AddRow(
				event1.ID, event1.Title, event1.Descr, event1.Owner, event1.StartAt, event1.EndAt, event1.SendNotificationAt,
			).AddRow(
				event2.ID, event2.Title, event2.Descr, event2.Owner, event2.StartAt, event2.EndAt, event2.SendNotificationAt,
			),
		)

	result, err := mockStorage.ListEventsOnARange(event1.StartAt, tMax)
	require.NoError(t, err)
	expected := storage.EventsSlice{event1, event2}
	require.Equal(t, expected, result)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestListMonthEvents(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	mockReturnRows := sqlmock.NewRows(eventColumns)
	expected := storage.EventsSlice{}

	for i := 0; i < 5; i++ {
		event := events[0]
		event.StartAt = event.StartAt.Add(time.Hour * 24 * 10 * time.Duration(i))
		if i < 30 {
			expected = append(expected, event)
		}
		mockReturnRows.AddRow(event.ID, event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt)
	}

	tMax := events[0].StartAt.Add(7 * time.Hour * 24)

	mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE start_at BETWEEN", EventsTable)).
		WithArgs(events[0].StartAt.Format(storage.RequestDateTimeFormat), tMax.Format(storage.RequestDateTimeFormat)).
		WillReturnRows(mockReturnRows)

	result, err := mockStorage.ListEventsOnARange(events[0].StartAt, tMax)
	require.NoError(t, err)
	require.Equal(t, expected, result)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestBusyTimeAdd(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	event := events[0]

	mock.ExpectExec(fmt.Sprintf("SELECT (.+) FROM %s", EventsTable)).WithArgs(
		event.StartAt,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err := mockStorage.Add(event)
	require.ErrorIs(t, err, storage.ErrDateBusy)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
