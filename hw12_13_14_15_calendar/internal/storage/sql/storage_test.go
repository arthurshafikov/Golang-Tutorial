package sqlstorage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

var events = []storage.Event{
	{
		ID:        "1",
		Title:     "test",
		Descr:     "testdes2",
		Owner:     42,
		StartDate: "2019-12-31",
		StartTime: "16:50:02.457276",
		EndDate:   "2019-12-31",
		EndTime:   "16:50:02.457276",
	},
	{
		ID:        "2",
		Title:     "test2",
		Descr:     "testdes2",
		Owner:     13,
		StartDate: "2020-12-31",
		StartTime: "10:00:00.457276",
		EndDate:   "2019-12-31",
		EndTime:   "16:50:02.457276",
	},
}

func TestSqlStorage(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mockStorage := Storage{
		dsn: "test dsn",
		db:  sqlxDB,
	}

	t.Run("test add", func(t *testing.T) {
		event := events[0]
		mock.ExpectExec("SELECT (.+) FROM events").WithArgs(
			event.StartDate, event.StartTime,
		).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("INSERT INTO events").WithArgs(
			event.Title, event.Descr, event.Owner,
			event.StartDate, event.EndDate, event.SendNotificationAt,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = mockStorage.Add(event)
		require.NoError(t, err)

		err := mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("test change", func(t *testing.T) {
		event := events[0]
		mock.ExpectExec("SELECT (.+) FROM events").WithArgs(
			event.StartDate, event.StartTime,
		).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE events SET").WithArgs(
			event.Title, event.Descr, event.Owner,
			event.StartDate, event.EndDate, event.SendNotificationAt,
			event.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = mockStorage.Change(event)
		require.NoError(t, err)

		err := mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("test get", func(t *testing.T) {
		event := events[0]
		columns := []string{
			"id", "title", "descr", "owner", "start_date", "start_time", "end_date", "end_time", "send_notification_at",
		}
		mock.ExpectQuery("SELECT (.+) FROM events").WithArgs(event.ID).WillReturnRows(
			sqlmock.NewRows(columns).AddRow(
				event.ID, event.Title, event.Descr, event.Owner,
				event.StartDate, event.StartTime, event.EndDate,
				event.EndTime, event.SendNotificationAt,
			),
		)
		resEvent, err := mockStorage.Get(event.ID)
		require.NoError(t, err)
		require.Equal(t, event, resEvent)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("test delete", func(t *testing.T) {
		event := events[0]
		mock.ExpectExec("DELETE FROM events").WithArgs(
			event.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = mockStorage.Delete(event)
		require.NoError(t, err)

		err := mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("test busy time add", func(t *testing.T) {
		event := events[0]
		mock.ExpectExec("SELECT (.+) FROM events").WithArgs(
			event.StartDate, event.StartTime,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = mockStorage.Add(event)
		require.ErrorIs(t, err, storage.ErrDateBusy)

		err := mock.ExpectationsWereMet()
		require.NoError(t, err)
	})
}
