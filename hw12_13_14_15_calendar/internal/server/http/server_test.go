package internalhttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/app"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/thewolf27/hw12_13_14_15_calendar/internal/storage/memory"
)

const expectedSuccessJSON = `{"data":"Success","error":null}`

var event = storage.Event{
	ID:      1,
	Title:   "test",
	Descr:   "testdes2",
	Owner:   42,
	StartAt: time.Date(2019, 04, 31, 05, 13, 58, 10, time.Time{}.UTC().Location()),
	EndAt:   time.Date(2019, 05, 06, 05, 13, 58, 10, time.Time{}.UTC().Location()),
}

var eventReq = storage.EventRequest{
	ID:                 1,
	Title:              "test",
	Descr:              "testdes2",
	Owner:              42,
	StartAt:            time.Date(2019, 04, 31, 05, 13, 58, 10, time.Time{}.UTC().Location()).Format(storage.RequestDateTimeFormat),
	EndAt:              time.Date(2019, 05, 06, 05, 13, 58, 10, time.Time{}.UTC().Location()).Format(storage.RequestDateTimeFormat),
	SendNotificationAt: time.Date(2019, 05, 06, 05, 13, 58, 10, time.Time{}.UTC().Location()).Format(storage.RequestDateTimeFormat),
}

type loggerMock struct{}

func (l loggerMock) Info(msg string)  {}
func (l loggerMock) Warn(msg string)  {}
func (l loggerMock) Error(msg string) {}

func NewServerMock() *Server {
	l := loggerMock{}
	m := memorystorage.New()
	app := app.New(l, m)
	return NewServer(l, app, "localhost", "9999")
}

func TestCreateEvent(t *testing.T) {
	server := NewServerMock()

	body, err := json.Marshal(eventReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.Create(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, expectedSuccessJSON, string(data))
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUpdateEvent(t *testing.T) {
	server := NewServerMock()

	server.App.CreateEvent(event)

	eventReq.ID = event.ID
	eventReq.Title = "New Title"
	body, err := json.Marshal(eventReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.Update(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, expectedSuccessJSON, string(data))
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func TestDeleteEvent(t *testing.T) {
	server := NewServerMock()

	server.App.CreateEvent(event)

	body, err := json.Marshal(eventReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/delete", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.Delete(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, expectedSuccessJSON, string(data))
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func TestListADayEvents(t *testing.T) {
	server := NewServerMock()

	server.App.CreateEvent(event)

	request := struct {
		Date string `json:"date"`
	}{
		Date: event.StartAt.Format(storage.RequestDateFormat),
	}

	body, err := json.Marshal(request)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/list-a-day", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.ListEventsOnADay(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	expected, err := json.Marshal(serverResponse{
		Data: storage.EventsSlice{event},
	})
	require.NoError(t, err)
	require.Equal(t, expected, data)
	require.Equal(t, http.StatusOK, res.StatusCode)
}
