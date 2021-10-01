package app

import (
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/config"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
	Config  config.Config
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
	GetEventsThatNeedToBeSend(time.Time) (storage.EventsSlice, error)
	GetEventsWhereEndAtBeforeGivenTimestamp(time.Time) (storage.EventsSlice, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(event storage.Event) error {
	return a.Storage.Add(event)
}

func (a *App) UpdateEvent(event storage.Event) error {
	return a.Storage.Change(event)
}

func (a *App) DeleteEvent(event storage.Event) error {
	return a.Storage.Delete(event)
}

func (a *App) ListEventsOnADay(date time.Time) (storage.EventsSlice, error) {
	return a.Storage.ListEventsOnADay(date)
}

func (a *App) ListEventsOnAWeek(timeStart time.Time) (storage.EventsSlice, error) {
	timePlusWeek := timeStart.Add(time.Hour * 24 * 7)
	return a.Storage.ListEventsOnARange(timeStart, timePlusWeek)
}

func (a *App) ListEventsOnAMonth(timeStart time.Time) (storage.EventsSlice, error) {
	timePlusMonth := timeStart.Add(time.Hour * 24 * 30)
	return a.Storage.ListEventsOnARange(timeStart, timePlusMonth)
}
