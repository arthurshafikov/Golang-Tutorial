package app

import (
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	Add(storage.Event) error
	Change(storage.Event) error
	Get(storage.Event) (storage.Event, error)
	Delete(storage.Event) error
	ListEventsOnADay(time.Time) (storage.EventsSlice, error)
	ListEventsOnARange(time.Time, time.Time) (storage.EventsSlice, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(event storage.Event) error { // ctx context.Context??
	return a.Storage.Add(event)
}

func (a *App) UpdateEvent(event storage.Event) error { // ctx context.Context??
	return a.Storage.Change(event)
}

func (a *App) DeleteEvent(event storage.Event) error { // ctx context.Context??
	return a.Storage.Delete(event)
}

func (a *App) ListEventsOnADay(t time.Time) (storage.EventsSlice, error) { // ctx context.Context??
	return a.Storage.ListEventsOnADay(t)
}

func (a *App) ListEventsOnAWeek(t time.Time) (storage.EventsSlice, error) { // ctx context.Context??
	tMax := t.Add(7 * time.Hour * 24)
	return a.Storage.ListEventsOnARange(t, tMax)
}

func (a *App) ListEventsOnAMonth(t time.Time) (storage.EventsSlice, error) { // ctx context.Context??
	tMax := t.Add(30 * time.Hour * 24)
	return a.Storage.ListEventsOnARange(t, tMax)
}
