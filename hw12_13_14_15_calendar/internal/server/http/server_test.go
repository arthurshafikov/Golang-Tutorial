package internalhttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arthurshafikov/hw12_13_14_15_calendar/internal/server"
	"github.com/arthurshafikov/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
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

var eventReq = server.EventRequest{
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

type appMock struct{}

func (app appMock) CreateEvent(event storage.Event) (int64, error) {
	return 0, nil
}

func (app appMock) UpdateEvent(event storage.Event) (int64, error) {
	return 0, nil
}

func (app appMock) DeleteEvent(event storage.Event) error {
	return nil
}

func (app appMock) ListEventsOnADay(date time.Time) (storage.EventsSlice, error) {
	return storage.EventsSlice{event}, nil
}

func (app appMock) ListEventsOnAWeek(startAt time.Time) (storage.EventsSlice, error) {
	return storage.EventsSlice{}, nil
}

func (app appMock) ListEventsOnAMonth(startAt time.Time) (storage.EventsSlice, error) {
	return storage.EventsSlice{}, nil
}

func NewServerMock() *Server {
	l := loggerMock{}
	app := appMock{}

	return NewServer(l, app, "localhost", "9999")
}

func TestCreateEvent(t *testing.T) {
	server := NewServerMock()

	body, err := json.Marshal(eventReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.create(w, req)

	res := w.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, expectedSuccessJSON, string(data))
}

func TestUpdateEvent(t *testing.T) {
	server := NewServerMock()

	body, err := json.Marshal(eventReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.update(w, req)

	res := w.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, expectedSuccessJSON, string(data))
}

func TestDeleteEvent(t *testing.T) {
	server := NewServerMock()

	body, err := json.Marshal(eventReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/delete", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.delete(w, req)

	res := w.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, expectedSuccessJSON, string(data))
}

func TestListADayEvents(t *testing.T) {
	server := NewServerMock()

	request := struct {
		Date string `json:"date"`
	}{
		Date: event.StartAt.Format(storage.RequestDateFormat),
	}

	body, err := json.Marshal(request)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/list-a-day", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.listEventsOnADay(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	expected, err := json.Marshal(serverResponse{
		Data: storage.EventsSlice{event},
	})
	require.NoError(t, err)
	require.Equal(t, expected, data)
}
