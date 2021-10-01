package memorystorage

import (
	"sync"
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	Events storage.EventsSlice
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Add(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.checkIfStartAtIsBusy(event) {
		return storage.ErrStartAtBusy
	}
	s.Events = append(s.Events, event)
	return nil
}

func (s *Storage) Change(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, e := range s.Events {
		if e.ID == event.ID {
			if s.checkIfStartAtIsBusy(event) {
				return storage.ErrStartAtBusy
			}
			s.Events[i] = event
		}
	}

	return nil
}

func (s *Storage) Get(event storage.Event) (storage.Event, error) {
	for _, e := range s.Events {
		if event.ID == e.ID {
			return e, nil
		}
	}

	return storage.Event{}, storage.ErrNotFound
}

func (s *Storage) Delete(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, e := range s.Events {
		if e.ID == event.ID {
			s.Events = append(s.Events[:i], s.Events[i+1:]...)
			return nil
		}
	}

	return storage.ErrNotFound
}

func (s *Storage) ListEventsOnADay(date time.Time) (storage.EventsSlice, error) {
	events := storage.EventsSlice{}
	for _, e := range s.Events {
		if date.Format(storage.RequestDateFormat) == e.StartAt.Format(storage.RequestDateFormat) {
			events = append(events, e)
		}
	}

	return events, nil
}

func (s *Storage) ListEventsOnARange(timeStart, timePlusRange time.Time) (storage.EventsSlice, error) {
	events := storage.EventsSlice{}
	for _, e := range s.Events {
		if (timePlusRange.After(e.StartAt) && timeStart.Before(e.StartAt)) ||
			timeStart.Equal(e.StartAt) {
			events = append(events, e)
		}
	}

	return events, nil
}

func (s *Storage) GetEventsThatNeedToBeSend(timeTo time.Time) (storage.EventsSlice, error) {
	events := storage.EventsSlice{}
	for _, e := range s.Events {
		if timeTo.After(e.SendNotificationAt) && !e.IsSent {
			events = append(events, e)
		}
	}

	return events, nil
}

func (s *Storage) GetEventsWhereEndAtBeforeGivenTimestamp(timeTo time.Time) (storage.EventsSlice, error) {
	events := storage.EventsSlice{}
	for _, e := range s.Events {
		if timeTo.After(e.EndAt) {
			events = append(events, e)
		}
	}

	return events, nil
}

func (s *Storage) checkIfStartAtIsBusy(event storage.Event) bool {
	for _, e := range s.Events {
		if e.StartAt == event.StartAt && e.ID != event.ID {
			return true
		}
	}

	return false
}
