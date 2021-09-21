package memorystorage

import (
	"fmt"
	"sync"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

type EventSlice []storage.Event

type Storage struct {
	mu     sync.RWMutex
	Events EventSlice
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
			if (event.StartDate != e.StartDate || event.StartTime != e.StartTime) && s.checkIfStartDateTimeIsBusy(event) {
				return storage.ErrDateBusy
			}
			s.Events[i] = event
		}
	}
	return nil
}

func (s *Storage) Get(id string) (storage.Event, error) {
	fmt.Println(s.Events)
	for _, e := range s.Events {
		if id == e.ID {
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
		}
	}
	return nil
}

func (s *Storage) List() error {
	for _, e := range s.Events {
		fmt.Println(e)
	}
	return nil
}

func (s *Storage) checkIfStartDateTimeIsBusy(event storage.Event) bool {
	for _, e := range s.Events {
		if e.StartDate == event.StartDate && e.StartTime == event.StartTime {
			return true
		}
	}
	return false
}
