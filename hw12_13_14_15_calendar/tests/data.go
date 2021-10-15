package tests

import (
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

var (
	eventTime            = time.Date(2018, 04, 31, 0, 0, 0, 0, time.FixedZone("", 0)).Round(time.Second)
	nowTime              = time.Now().In(time.FixedZone("", 0))
	nowTimeMinusTwoYears = nowTime.Add(-1 * time.Hour * 24 * 730).Round(time.Second)
)

var event = storage.Event{
	ID:                 1,
	Title:              "Integration event",
	Descr:              "Integration event",
	Owner:              10,
	StartAt:            eventTime,
	EndAt:              eventTime.Add(time.Hour * 24).Round(time.Second),
	SendNotificationAt: eventTime,
}

var eventSlice = storage.EventsSlice{
	storage.Event{
		Title:              "Should be deleted from database event 1",
		Descr:              "Should be deleted from database event 1",
		Owner:              10,
		StartAt:            nowTimeMinusTwoYears,
		EndAt:              nowTimeMinusTwoYears,
		SendNotificationAt: nowTimeMinusTwoYears,
	},
	storage.Event{
		Title:              "Should be deleted from database event 2",
		Descr:              "Should be deleted from database event 2",
		Owner:              10,
		StartAt:            nowTimeMinusTwoYears.Add(5 * time.Hour).Round(time.Second),
		EndAt:              nowTimeMinusTwoYears.Add(5 * time.Hour * 2).Round(time.Second),
		SendNotificationAt: nowTimeMinusTwoYears.Add(5 * time.Hour).Round(time.Second),
	},
	storage.Event{
		Title:              "Should stay in database event 1",
		Descr:              "Should stay in database event 1",
		Owner:              10,
		StartAt:            nowTime.Add(-1 * time.Hour * 24 * 30).Round(time.Second),
		EndAt:              nowTime.Add(-1 * time.Hour * 24 * 30).Add(5 * time.Minute).Round(time.Second),
		SendNotificationAt: nowTime.Add(-1 * time.Hour * 24 * 30).Round(time.Second),
	},
}
