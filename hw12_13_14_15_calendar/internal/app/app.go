package app

import (
	"context"

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
	Get(id string) (storage.Event, error)
	Delete(storage.Event) error
	List() error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	/*
		event := abstractstorage.Event{
			ID:         "1",
			Title:      "test",
			Descr:      "testdes",
			Owner:      42,
			Start_date: "2019-12-31",
			Start_time: "16:50:02.457276",
			End_date:   "2019-12-31",
			End_time:   "16:50:02.457276",
		}
		storage.(*sqlstorage.Storage).Connect(context.Background())
		storage.Add(event)
	*/
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
