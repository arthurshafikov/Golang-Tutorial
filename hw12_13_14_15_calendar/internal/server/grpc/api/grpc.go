package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	grpcapi "github.com/thewolf27/hw12_13_14_15_calendar/internal/server/grpc/generated"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

var (
	EmptyErrorResponse = "null"
	successResponse    = &grpcapi.ServerResponse{
		Data:  storage.RequestSuccessMessage,
		Error: EmptyErrorResponse,
	}
)

type CalendarService struct {
	app Application
	grpcapi.UnimplementedCalendarServiceServer
}

func (c *CalendarService) Create(ctx context.Context, req *grpcapi.EventRequest) (*grpcapi.ServerResponse, error) {
	event, err := storage.ParseEvent(req.Event, true)
	if err != nil {
		return nil, err
	}
	if err := c.app.CreateEvent(event); err != nil {
		return nil, fmt.Errorf("didnt created the event %w", err)
	}
	return successResponse, nil
}

func (c *CalendarService) Update(ctx context.Context, req *grpcapi.EventRequest) (*grpcapi.ServerResponse, error) {
	event, err := storage.ParseEvent(req.Event, false)
	if err != nil {
		return nil, err
	}
	if err := c.app.UpdateEvent(event); err != nil {
		return nil, fmt.Errorf("didnt updated the event %w", err)
	}
	return successResponse, nil
}

func (c *CalendarService) Delete(ctx context.Context, req *grpcapi.EventRequest) (*grpcapi.ServerResponse, error) {
	event, err := storage.ParseEvent(req.Event, false)
	if err != nil {
		return nil, err
	}
	if err := c.app.DeleteEvent(event); err != nil {
		return nil, fmt.Errorf("didnt deleted the event %w", err)
	}
	return successResponse, nil
}

func (c *CalendarService) ListEventsOnADay(ctx context.Context, req *grpcapi.ListEventsRequest) (*grpcapi.ServerResponse, error) {
	return c.listEvents(req, c.app.ListEventsOnADay)
}

func (c *CalendarService) ListEventsOnAWeek(ctx context.Context, req *grpcapi.ListEventsRequest) (*grpcapi.ServerResponse, error) {
	return c.listEvents(req, c.app.ListEventsOnAWeek)
}

func (c *CalendarService) ListEventsOnAMonth(ctx context.Context, req *grpcapi.ListEventsRequest) (*grpcapi.ServerResponse, error) {
	return c.listEvents(req, c.app.ListEventsOnAMonth)
}

func (c *CalendarService) listEvents(req *grpcapi.ListEventsRequest, listEventsFunc storage.ListEventsFunction) (*grpcapi.ServerResponse, error) {
	date, err := time.Parse(storage.RequestDateFormat, req.Date)
	if err != nil {
		return nil, fmt.Errorf("didnt parsed the date %w", err)
	}
	eventsSlice, err := listEventsFunc(date)
	if err != nil {
		return nil, fmt.Errorf("didnt listed the week events %w", err)
	}
	events, err := json.Marshal(eventsSlice)
	if err != nil {
		return nil, fmt.Errorf("didnt marshal the events %w", err)
	}
	return &grpcapi.ServerResponse{
		Data:  string(events),
		Error: EmptyErrorResponse,
	}, nil
}
