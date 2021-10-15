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
	"id", "title", "descr", "owner", "start_at", "end_at", "send_notification_at", "is_sent",
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
		event.StartAt, 0,
	).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", EventsTable)).WithArgs(
		event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt, EventsIsSentFalse,
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err := mockStorage.Add(event)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestChangeEvent(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	event := events[0]

	mock.ExpectExec(fmt.Sprintf("SELECT (.+) FROM %s", EventsTable)).WithArgs(
		event.StartAt, event.ID,
	).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(fmt.Sprintf("UPDATE %s SET", EventsTable)).WithArgs(
		event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt, event.IsSent, event.ID,
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err := mockStorage.Change(event)
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
			event.ID, event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt, event.IsSent,
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

func TestListEventsOnADay(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	dateToListAt := time.Date(2019, 04, 10, 0, 0, 0, 0, time.Time{}.UTC().Location())

	event1 := events[0]
	event1.StartAt = dateToListAt
	event2 := events[0]
	event2.StartAt = dateToListAt.Add(24 * time.Hour)

	mock.ExpectQuery(
		fmt.Sprintf("SELECT (.+) FROM %s WHERE DATE\\(start_at\\)", EventsTable),
	).WithArgs(dateToListAt.Format(storage.RequestDateFormat)).WillReturnRows(
		sqlmock.NewRows(eventColumns).AddRow(
			event1.ID, event1.Title, event1.Descr, event1.Owner, event1.StartAt, event1.EndAt, event1.SendNotificationAt, event1.IsSent,
		).AddRow(
			event2.ID, event2.Title, event2.Descr, event2.Owner, event2.StartAt, event2.EndAt, event2.SendNotificationAt, event2.IsSent,
		),
	)

	result, err := mockStorage.ListEventsOnADay(event1.StartAt)
	require.NoError(t, err)
	expected := storage.EventsSlice{event1, event2}
	require.Equal(t, expected, result)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestListEventsOnARange(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	expectedEvents := storage.EventsSlice{}

	rangeDifferenceInDays := 3
	rangeStartTime := time.Date(2020, 01, 10, 0, 0, 0, 0, time.Now().Location())
	rangeEndTime := rangeStartTime.Add(time.Duration(rangeDifferenceInDays) * time.Hour * 24)

	returnResultSQL := sqlmock.NewRows(eventColumns)
	for i := 0; i < 10; i++ {
		event := events[0]
		event.StartAt = rangeStartTime.Add(time.Duration(i) * time.Hour * 24)

		if i < rangeDifferenceInDays {
			returnResultSQL.AddRow(
				event.ID, event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt, event.IsSent,
			)
			expectedEvents = append(expectedEvents, event)
		}
	}

	mock.ExpectQuery(
		fmt.Sprintf("SELECT (.+) FROM %s WHERE start_at BETWEEN", EventsTable),
	).WithArgs(rangeStartTime.Format(storage.RequestDateTimeFormat), rangeEndTime.Format(storage.RequestDateTimeFormat)).WillReturnRows(
		returnResultSQL,
	)

	result, err := mockStorage.ListEventsOnARange(rangeStartTime, rangeEndTime)
	require.NoError(t, err)
	require.Equal(t, expectedEvents, result)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetEventsThatNeedToBeSend(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	mockReturnRows := sqlmock.NewRows(eventColumns)
	expected := storage.EventsSlice{}

	for i := 0; i < 9; i++ {
		event := events[0]
		event.StartAt = event.StartAt.Add(time.Hour * 24 * 10 * time.Duration(i))
		event.SendNotificationAt = time.Now().Add(time.Duration(i-4) * time.Hour * 24)
		if i < 5 {
			mockReturnRows.AddRow(event.ID, event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt, event.IsSent)
			expected = append(expected, event)
		}
	}

	mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE send_notification_at <=", EventsTable)).
		WithArgs(time.Now().Format(storage.RequestDateTimeFormat)).
		WillReturnRows(mockReturnRows)

	result, err := mockStorage.GetEventsThatNeedToBeSend(time.Now())
	require.NoError(t, err)
	require.Equal(t, expected, result)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetEventsWhereEndAtBeforeGivenTimestamp(t *testing.T) {
	mockDB, mock, mockStorage := newSQLStorageMock(t)
	defer mockDB.Close()

	mockReturnRows := sqlmock.NewRows(eventColumns)
	expected := storage.EventsSlice{}

	for i := 0; i < 9; i++ {
		event := events[0]
		event.StartAt = event.StartAt.Add(time.Hour * 24 * 10 * time.Duration(i))
		event.EndAt = time.Now().Add(time.Duration(i-4) * time.Hour * 24)
		if i < 5 {
			mockReturnRows.AddRow(event.ID, event.Title, event.Descr, event.Owner, event.StartAt, event.EndAt, event.SendNotificationAt, event.IsSent)
			expected = append(expected, event)
		}
	}

	mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE end_at <=", EventsTable)).
		WithArgs(time.Now().Format(storage.RequestDateTimeFormat)).
		WillReturnRows(mockReturnRows)

	result, err := mockStorage.GetEventsWhereEndAtBeforeGivenTimestamp(time.Now())
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
		event.StartAt, 0,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := mockStorage.Add(event)
	require.ErrorIs(t, err, storage.ErrStartAtBusy)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
