package tests

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	sender "github.com/thewolf27/hw12_13_14_15_calendar/internal/app"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/thewolf27/hw12_13_14_15_calendar/internal/storage/sql"
)

type appSuiteHandler struct {
	AppSuite
}

func TestAppSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test app")
	}

	appSuiteHandler := &appSuiteHandler{
		AppSuite: AppSuite{},
	}

	suite.Run(t, appSuiteHandler)
}

func (app *appSuiteHandler) TestAddEventsScheduleItSendAndDeleteOld() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	go app.Scheduler.Run(ctx)
	go app.Sender.Run(ctx)

	for i, e := range eventSlice {
		eventID, err := app.App.CreateEvent(e)
		require.NoError(app.T(), err)
		eventSlice[i].ID = eventID
		eventSlice[i].IsSent = true
	}

	messagesCounter := 0
	for message := range app.Sender.ConsumerMessagesCh {
		e := eventSlice[messagesCounter]
		notification := sender.Notification{
			ID:    e.ID,
			Title: e.Title,
			Owner: e.Owner,
			Date:  e.StartAt.Format(storage.RequestDateTimeFormat),
		}

		expectedJSONString, err := json.Marshal(notification)
		require.NoError(app.T(), err)

		require.Equal(app.T(), string(expectedJSONString), message)
		messagesCounter++

		if messagesCounter == len(eventSlice) {
			cancel()
			time.Sleep(time.Second * 2) // sleep while events are deleting
		}
	}

	expectedMissingEvents := eventSlice[:1]
	for _, e := range expectedMissingEvents {
		_, err := app.App.Storage.Get(e)
		require.ErrorIs(app.T(), err, storage.ErrNotFound)
	}

	expectedStayEvent := eventSlice[2]

	databaseEvent, err := app.App.Storage.Get(expectedStayEvent)
	require.NoError(app.T(), err)
	require.Equal(app.T(), expectedStayEvent, databaseEvent)

	cancel()
	err = app.App.Storage.(*sqlstorage.Storage).ResetDB()
	require.NoError(app.T(), err)
}
