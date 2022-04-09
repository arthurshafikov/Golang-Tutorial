package tests

import (
	"testing"
	"time"

	"github.com/arthurshafikov/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/arthurshafikov/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type calendarSuiteHandler struct {
	CalendarSuite
}

func TestCalendarSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test app")
	}

	appSuiteHandler := &calendarSuiteHandler{
		CalendarSuite: CalendarSuite{},
	}

	suite.Run(t, appSuiteHandler)
}

func (cal *calendarSuiteHandler) TestAddEvent() {
	id, err := cal.App.CreateEvent(event)
	require.NoError(cal.T(), err)
	event.ID = id

	databaseEvent, err := cal.App.Storage.Get(storage.Event{ID: id})
	require.NoError(cal.T(), err)
	require.Equal(cal.T(), event, databaseEvent)

	err = cal.App.Storage.(*sqlstorage.Storage).ResetDB()
	require.NoError(cal.T(), err)
}

func (cal *calendarSuiteHandler) TestUpdateEvent() {
	id, err := cal.App.CreateEvent(event)
	require.NoError(cal.T(), err)

	newTitle := "New Title"
	newEvent := event
	newEvent.Title = newTitle
	newEvent.ID = id
	id, err = cal.App.UpdateEvent(newEvent)
	require.NoError(cal.T(), err)

	databaseEvent, err := cal.App.Storage.Get(storage.Event{ID: id})
	require.NoError(cal.T(), err)
	require.Equal(cal.T(), newTitle, databaseEvent.Title)

	err = cal.App.Storage.(*sqlstorage.Storage).ResetDB()
	require.NoError(cal.T(), err)
}

func (cal *calendarSuiteHandler) TestListADayEvents() {
	listEvent := event
	dateToListAt := time.Date(2020, 07, 15, 22, 0, 0, 0, time.FixedZone("", 0))
	eventsCount := 10
	expectedEvents := storage.EventsSlice{}

	for i := 0; i < eventsCount; i++ {
		listEvent.StartAt = dateToListAt.Add(time.Hour * time.Duration(i))
		id, err := cal.App.CreateEvent(listEvent)
		require.NoError(cal.T(), err)
		listEvent.ID = id
		if listEvent.StartAt.Day() == dateToListAt.Day() {
			expectedEvents = append(expectedEvents, listEvent)
		}
	}

	events, err := cal.App.ListEventsOnADay(dateToListAt)
	require.NoError(cal.T(), err)
	require.Equal(cal.T(), expectedEvents, events)

	err = cal.App.Storage.(*sqlstorage.Storage).ResetDB()
	require.NoError(cal.T(), err)
}

func (cal *calendarSuiteHandler) TestListAWeekEvents() {
	listEvent := event
	startOfTheWeek := time.Date(2020, 10, 10, 22, 0, 0, 0, time.FixedZone("", 0))
	eventsCount := 14
	daysInWeek := 7
	expectedEvents := storage.EventsSlice{}

	for i := 0; i < eventsCount; i++ {
		listEvent.StartAt = startOfTheWeek.Add(time.Duration(i) * time.Hour * 24)
		id, err := cal.App.CreateEvent(listEvent)
		require.NoError(cal.T(), err)
		listEvent.ID = id
		if i <= daysInWeek {
			expectedEvents = append(expectedEvents, listEvent)
		}
	}

	events, err := cal.App.ListEventsOnAWeek(startOfTheWeek)
	require.NoError(cal.T(), err)
	require.Equal(cal.T(), expectedEvents, events)

	err = cal.App.Storage.(*sqlstorage.Storage).ResetDB()
	require.NoError(cal.T(), err)
}

func (cal *calendarSuiteHandler) TestListAMonthEvents() {
	listEvent := event
	startOfTheMonth := time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("", 0))
	expectedEvents := storage.EventsSlice{}
	eventsCount := 14
	daysInMonth := 31
	spreadValue := 8

	for i := 0; i < eventsCount; i++ {
		listEvent.StartAt = startOfTheMonth.Add(time.Duration(i*spreadValue) * time.Hour * 24)
		id, err := cal.App.CreateEvent(listEvent)
		require.NoError(cal.T(), err)

		listEvent.ID = id
		if i*spreadValue < daysInMonth {
			expectedEvents = append(expectedEvents, listEvent)
		}
	}

	events, err := cal.App.ListEventsOnAMonth(startOfTheMonth)
	require.NoError(cal.T(), err)
	require.Equal(cal.T(), expectedEvents, events)

	err = cal.App.Storage.(*sqlstorage.Storage).ResetDB()
	require.NoError(cal.T(), err)
}
