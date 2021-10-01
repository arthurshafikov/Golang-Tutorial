package memorystorage

import (
	"sync"
	"testing"
	"time"

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

func TestMemoryAddEvent(t *testing.T) {
	memstorage := New()

	for _, e := range events {
		err := memstorage.Add(e)
		require.NoError(t, err)
	}

	first, err := memstorage.Get(storage.Event{ID: 1})
	require.NoError(t, err)
	second, err := memstorage.Get(storage.Event{ID: 2})
	require.NoError(t, err)

	require.Equal(t, events[0], first)
	require.Equal(t, events[1], second)
}

func TestMemoryChangeEvent(t *testing.T) {
	memstorage := New()

	event := events[0]

	err := memstorage.Add(event)
	require.NoError(t, err)

	event.Title = "New Title"
	err = memstorage.Change(event)
	require.NoError(t, err)

	first, err := memstorage.Get(event)
	require.NoError(t, err)

	require.Equal(t, "New Title", first.Title)
}

func TestMemoryDeleteEvent(t *testing.T) {
	memstorage := New()
	for _, e := range events {
		err := memstorage.Add(e)
		require.NoError(t, err)
	}

	for _, e := range events {
		err := memstorage.Delete(e)
		require.NoError(t, err)
	}

	first, err := memstorage.Get(storage.Event{ID: 1})
	require.ErrorIs(t, err, storage.ErrNotFound)
	second, err := memstorage.Get(storage.Event{ID: 2})
	require.ErrorIs(t, err, storage.ErrNotFound)

	require.Equal(t, storage.Event{}, first)
	require.Equal(t, storage.Event{}, second)
}

func TestMemoryListDayEvents(t *testing.T) {
	memstorage := New()

	for _, e := range events {
		err := memstorage.Add(e)
		require.NoError(t, err)
	}

	event := events[0]
	events := storage.EventsSlice{event}

	result, err := memstorage.ListEventsOnADay(event.StartAt)
	require.NoError(t, err)
	require.Equal(t, events, result)
}

func TestMemoryListWeekEvents(t *testing.T) {
	memstorage := New()
	weekEvents := events

	for i := 0; i < 2; i++ {
		e := &weekEvents[i]
		e.StartAt = time.Date(2020, 04, 31, 05, 13, 58, 10, time.Now().Location()).Add(time.Duration(i) * time.Hour * 24)
		err := memstorage.Add(*e)
		require.NoError(t, err)
	}

	eventTime := weekEvents[0].StartAt

	timePlusWeek := eventTime.Add(7 * time.Hour * 24)
	result, err := memstorage.ListEventsOnARange(eventTime, timePlusWeek)
	require.NoError(t, err)

	require.Equal(t, weekEvents, result)
}

func TestMemoryListMonthEvents(t *testing.T) {
	memstorage := New()
	expectEvents := storage.EventsSlice{}

	for i := 0; i < 9; i++ {
		event := events[0]
		if i < 5 {
			event.StartAt = event.StartAt.Add(time.Duration(i) * time.Hour * 24)
			expectEvents = append(expectEvents, event)
		} else {
			event.StartAt = event.StartAt.AddDate(i*10, 0, 0)
		}
		err := memstorage.Add(event)
		require.NoError(t, err)
	}

	eventTime := events[0].StartAt

	timePlusMonth := eventTime.Add(30 * time.Hour * 24)
	result, err := memstorage.ListEventsOnARange(eventTime, timePlusMonth)
	require.NoError(t, err)

	require.Equal(t, expectEvents, result)
}

func TestGetEventsThatNeedToBeSend(t *testing.T) {
	memstorage := New()
	expectEvents := storage.EventsSlice{}

	for i := 0; i < 9; i++ {
		e := events[0]
		e.SendNotificationAt = time.Now().Add(time.Duration(i-4) * time.Hour * 24)
		e.StartAt = e.StartAt.Add(time.Duration(i) * time.Hour * 24)
		err := memstorage.Add(e)
		require.NoError(t, err)
		if i < 5 {
			expectEvents = append(expectEvents, e)
		}
	}

	result, err := memstorage.GetEventsThatNeedToBeSend(time.Now())
	require.NoError(t, err)
	require.Equal(t, expectEvents, result)
}

func TestGetEventsWhereEndAtBeforeGivenTimestamp(t *testing.T) {
	memstorage := New()
	expectEvents := storage.EventsSlice{}

	for i := 0; i < 9; i++ {
		e := events[0]
		e.StartAt = time.Now().Add(time.Duration(i-4) * time.Hour * 24)
		e.EndAt = time.Now().Add(time.Duration(i-4) * time.Hour * 24)
		err := memstorage.Add(e)
		require.NoError(t, err)
		if i < 5 {
			expectEvents = append(expectEvents, e)
		}
	}

	result, err := memstorage.GetEventsWhereEndAtBeforeGivenTimestamp(time.Now())
	require.NoError(t, err)
	require.Equal(t, expectEvents, result)
}

func TestMemoryBusyTimeAddError(t *testing.T) {
	memstorage := New()
	for _, e := range events {
		err := memstorage.Add(e)
		require.NoError(t, err)
	}

	busyevent := events[0]
	busyevent.ID = 9999
	err := memstorage.Add(busyevent)
	require.ErrorIs(t, err, storage.ErrStartAtBusy)
}

func TestMemoryConcurrent(t *testing.T) {
	memstorage := New()
	var wg sync.WaitGroup
	wg.Add(2)

	for k := 0; k < 2; k++ {
		go func(k int) {
			defer wg.Done()
			event := events[k]
			for i := 0; i < 1000; i++ {
				event.StartAt = time.Now().AddDate(k*10, 0, i)
				err := memstorage.Add(event)
				require.NoError(t, err)
			}
		}(k)
	}

	wg.Wait()
	require.Equal(t, 2000, len(memstorage.Events))
}
