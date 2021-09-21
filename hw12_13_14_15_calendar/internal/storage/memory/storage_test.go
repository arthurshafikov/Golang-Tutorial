package memorystorage

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

var events = EventSlice{
	storage.Event{
		ID:        "1",
		Title:     "test",
		Descr:     "testdes2",
		Owner:     42,
		StartDate: "2019-12-31",
		StartTime: "16:50:02.457276",
		EndDate:   "2019-12-31",
		EndTime:   "16:50:02.457276",
	},
	storage.Event{
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

func TestMemoryStorage(t *testing.T) {
	t.Run("test add", func(t *testing.T) {
		memstorage := New()

		for _, e := range events {
			err := memstorage.Add(e)
			require.NoError(t, err)
		}

		first, err := memstorage.Get("1")
		require.NoError(t, err)
		second, err := memstorage.Get("2")
		require.NoError(t, err)

		require.Equal(t, events[0], first)
		require.Equal(t, events[1], second)
	})

	t.Run("test change", func(t *testing.T) {
		memstorage := New()

		event := events[0]

		err := memstorage.Add(event)
		require.NoError(t, err)
		event.Title = "New Title"
		err = memstorage.Change(event)
		require.NoError(t, err)

		first, err := memstorage.Get(event.ID)
		require.NoError(t, err)

		require.Equal(t, "New Title", first.Title)
	})

	t.Run("test Delete", func(t *testing.T) {
		memstorage := New()
		for _, e := range events {
			err := memstorage.Add(e)
			require.NoError(t, err)
		}

		for _, e := range events {
			err := memstorage.Delete(e)
			require.NoError(t, err)
		}

		first, err := memstorage.Get("1")
		require.ErrorIs(t, err, storage.ErrNotFound)
		second, err := memstorage.Get("2")
		require.ErrorIs(t, err, storage.ErrNotFound)

		require.Equal(t, storage.Event{}, first)
		require.Equal(t, storage.Event{}, second)
	})

	t.Run("test busy time", func(t *testing.T) {
		memstorage := New()
		for _, e := range events {
			err := memstorage.Add(e)
			require.NoError(t, err)
		}

		busyevent := events[0]
		err := memstorage.Add(busyevent)
		require.ErrorIs(t, err, storage.ErrDateBusy)
	})

	t.Run("test concurrent", func(t *testing.T) {
		memstorage := New()
		var wg sync.WaitGroup
		wg.Add(2)

		for k := 0; k < 2; k++ {
			go func(k int) {
				defer wg.Done()
				event := events[k]
				for i := 0; i < 1000; i++ {
					event.StartDate = time.Now().AddDate(k*10, 0, i).Format("2006-01-02")
					err := memstorage.Add(event)
					require.NoError(t, err)
				}
			}(k)
		}

		wg.Wait()
		require.Equal(t, 2000, len(memstorage.Events))
	})
}
