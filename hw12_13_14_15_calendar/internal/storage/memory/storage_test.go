package memorystorage

import (
	"sync"
	"testing"
	"time"

	"github.com/arthurshafikov/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

var eventsSlice = storage.EventsSlice{
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

	for _, e := range eventsSlice {
		_, err := memstorage.Add(e)
		require.NoError(t, err)
	}

	first, err := memstorage.Get(storage.Event{ID: eventsSlice[0].ID})
	require.NoError(t, err)
	require.Equal(t, eventsSlice[0], first)

	second, err := memstorage.Get(storage.Event{ID: eventsSlice[1].ID})
	require.NoError(t, err)
	require.Equal(t, eventsSlice[1], second)
}

func TestMemoryChangeEvent(t *testing.T) {
	memstorage := New()

	event := eventsSlice[0]

	id, err := memstorage.Add(event)
	require.NoError(t, err)

	event.Title = "New Title"
	_, err = memstorage.Change(event)
	require.NoError(t, err)

	first, err := memstorage.Get(storage.Event{ID: id})
	require.NoError(t, err)

	require.Equal(t, "New Title", first.Title)
}

func TestMemoryDeleteEvent(t *testing.T) {
	memstorage := New()
	for _, e := range eventsSlice {
		_, err := memstorage.Add(e)
		require.NoError(t, err)
	}

	for _, e := range eventsSlice {
		err := memstorage.Delete(e)
		require.NoError(t, err)
	}

	first, err := memstorage.Get(storage.Event{ID: 1})
	require.ErrorIs(t, err, storage.ErrNotFound)
	require.Equal(t, storage.Event{}, first)

	second, err := memstorage.Get(storage.Event{ID: 2})
	require.ErrorIs(t, err, storage.ErrNotFound)
	require.Equal(t, storage.Event{}, second)
}

func TestMemoryListEventsOnADay(t *testing.T) {
	memstorage := New()
	expectedEvents := storage.EventsSlice{}
	dateToListAt := time.Date(2019, 04, 10, 0, 0, 0, 0, time.Time{}.UTC().Location())

	for i, e := range eventsSlice {
		e.StartAt = dateToListAt.Add(time.Duration(i) * time.Minute)
		expectedEvents = append(expectedEvents, e)
		_, err := memstorage.Add(e)
		require.NoError(t, err)
	}

	result, err := memstorage.ListEventsOnADay(dateToListAt)
	require.NoError(t, err)
	require.Equal(t, expectedEvents, result)
}

func TestMemoryListEventsOnARange(t *testing.T) {
	memstorage := New()
	expectEvents := storage.EventsSlice{}

	rangeDifferenceInDays := 5
	rangeStartTime := time.Date(2020, 01, 10, 0, 0, 0, 0, time.Now().Location())
	rangeEndTime := rangeStartTime.Add(time.Duration(rangeDifferenceInDays) * time.Hour * 24)

	event := eventsSlice[0]

	for i := 0; i < 15; i++ {
		event.StartAt = rangeStartTime.Add(time.Duration(i) * time.Hour * 24)
		_, err := memstorage.Add(event)
		require.NoError(t, err)
		if i < rangeDifferenceInDays {
			expectEvents = append(expectEvents, event)
		}
	}

	result, err := memstorage.ListEventsOnARange(rangeStartTime, rangeEndTime)
	require.NoError(t, err)
	require.Equal(t, expectEvents, result)
}

func TestGetEventsThatNeedToBeSend(t *testing.T) {
	memstorage := New()
	expectEvents := storage.EventsSlice{}

	needToBeSendEventsCount := 5
	for i := 1; i < 9; i++ {
		event := eventsSlice[0]
		event.SendNotificationAt = time.Now().Add(time.Duration(i-needToBeSendEventsCount) * time.Hour * 24)
		event.StartAt = event.StartAt.Add(time.Duration(i) * time.Hour * 24)
		_, err := memstorage.Add(event)
		require.NoError(t, err)
		if i <= needToBeSendEventsCount {
			expectEvents = append(expectEvents, event)
		}
	}

	result, err := memstorage.GetEventsThatNeedToBeSend(time.Now())
	require.NoError(t, err)
	require.Equal(t, expectEvents, result)
}

func TestGetEventsWhereEndAtBeforeGivenTimestamp(t *testing.T) {
	memstorage := New()
	expectEvents := storage.EventsSlice{}

	expectedEventsCount := 5
	for i := 1; i < 9; i++ {
		e := eventsSlice[0]
		e.StartAt = time.Now().Add(time.Duration(i-expectedEventsCount) * time.Hour * 24)
		e.EndAt = time.Now().Add(time.Duration(i-expectedEventsCount) * time.Hour * 24)
		_, err := memstorage.Add(e)
		require.NoError(t, err)
		if i <= expectedEventsCount {
			expectEvents = append(expectEvents, e)
		}
	}

	result, err := memstorage.GetEventsWhereEndAtBeforeGivenTimestamp(time.Now())
	require.NoError(t, err)
	require.Equal(t, expectEvents, result)
	require.Equal(t, expectedEventsCount, len(result))
}

func TestMemoryBusyTimeAddError(t *testing.T) {
	memstorage := New()
	for _, e := range eventsSlice {
		_, err := memstorage.Add(e)
		require.NoError(t, err)
	}

	busyevent := eventsSlice[0]
	busyevent.ID = 9999
	id, err := memstorage.Add(busyevent)
	require.ErrorIs(t, err, storage.ErrStartAtBusy)
	require.Equal(t, int64(0), id)
}

func TestMemoryConcurrent(t *testing.T) {
	memstorage := New()
	var wg sync.WaitGroup
	wg.Add(2)

	for k := 0; k < 2; k++ {
		go func(k int) {
			defer wg.Done()
			event := eventsSlice[k]
			for i := 0; i < 1000; i++ {
				event.StartAt = time.Now().AddDate(k*10, 0, i)
				_, err := memstorage.Add(event)
				require.NoError(t, err)
			}
		}(k)
	}

	wg.Wait()
	require.Equal(t, 2000, len(memstorage.Events))
}
