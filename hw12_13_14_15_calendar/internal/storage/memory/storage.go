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

	if s.checkIfStartDateTimeIsBusy(event) {
		return storage.ErrDateBusy
	}
	s.Events = append(s.Events, event)
	return nil
}

func (s *Storage) Change(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, e := range s.Events {
		if e.ID == event.ID {
			if !event.StartAt.Equal(e.StartAt) && s.checkIfStartDateTimeIsBusy(event) {
				return storage.ErrDateBusy
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

func (s *Storage) ListEventsOnADay(t time.Time) (storage.EventsSlice, error) {
	events := storage.EventsSlice{}
	for _, e := range s.Events {
		if t.Format(storage.RequestDateFormat) == e.StartAt.Format(storage.RequestDateFormat) {
			events = append(events, e)
		}
	}

	return events, nil
}

func (s *Storage) ListEventsOnARange(t time.Time, tMax time.Time) (storage.EventsSlice, error) {
	t = t.Add(-1)

	events := storage.EventsSlice{}
	for _, e := range s.Events {
		if tMax.After(e.StartAt) && t.Before(e.StartAt) {
			events = append(events, e)
		}
	}

	return events, nil
}

func (s *Storage) checkIfStartDateTimeIsBusy(event storage.Event) bool {
	for _, e := range s.Events {
		if e.StartAt == event.StartAt {
			return true
		}
	}

	return false
}
