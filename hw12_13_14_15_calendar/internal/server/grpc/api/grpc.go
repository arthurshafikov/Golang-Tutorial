package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/server"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/server/grpc/generated"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

var successResponse = &generated.ServerResponse{
	Data:  storage.RequestSuccessMessage,
	Error: server.EmptyErrorResponse,
}

type CalendarService struct {
	app    Application
	logger Logger
	generated.UnimplementedCalendarServiceServer
}

func (c *CalendarService) Create(ctx context.Context, req *generated.EventRequest) (*generated.ServerResponse, error) {
	event, err := server.ParseEvent(req.Event, true)
	if err != nil {
		return nil, err
	}
	if err := c.app.CreateEvent(event); err != nil {
		return nil, fmt.Errorf(server.ErrCantCreateEventFormat, err)
	}
	return successResponse, nil
}

func (c *CalendarService) Update(ctx context.Context, req *generated.EventRequest) (*generated.ServerResponse, error) {
	event, err := server.ParseEvent(req.Event, false)
	if err != nil {
		return nil, err
	}
	if err := c.app.UpdateEvent(event); err != nil {
		return nil, fmt.Errorf(server.ErrCantUpdateEventFormat, err)
	}
	return successResponse, nil
}

func (c *CalendarService) Delete(ctx context.Context, req *generated.EventRequest) (*generated.ServerResponse, error) {
	event, err := server.ParseEvent(req.Event, false)
	if err != nil {
		return nil, err
	}
	if err := c.app.DeleteEvent(event); err != nil {
		return nil, fmt.Errorf(server.ErrCantDeleteEventFormat, err)
	}
	return successResponse, nil
}

func (c *CalendarService) ListEventsOnADay(ctx context.Context, req *generated.ListEventsRequest) (*generated.ServerResponse, error) {
	return c.listEvents(req, c.app.ListEventsOnADay)
}

func (c *CalendarService) ListEventsOnAWeek(ctx context.Context, req *generated.ListEventsRequest) (*generated.ServerResponse, error) {
	return c.listEvents(req, c.app.ListEventsOnAWeek)
}

func (c *CalendarService) ListEventsOnAMonth(ctx context.Context, req *generated.ListEventsRequest) (*generated.ServerResponse, error) {
	return c.listEvents(req, c.app.ListEventsOnAMonth)
}

func (c *CalendarService) listEvents(req *generated.ListEventsRequest, listEventsFunc storage.ListEventsFunction) (*generated.ServerResponse, error) {
	requestDate, err := time.Parse(storage.RequestDateFormat, req.Date)
	if err != nil {
		return nil, fmt.Errorf(server.ErrCantParseDateFormat, err)
	}

	eventsSlice, err := listEventsFunc(requestDate)
	if err != nil {
		return nil, fmt.Errorf(server.ErrCantListEventsFormat, err)
	}

	events, err := json.Marshal(eventsSlice)
	if err != nil {
		return nil, fmt.Errorf(server.ErrCantMarshalEventsFormat, err)
	}

	return &generated.ServerResponse{
		Data:  string(events),
		Error: server.EmptyErrorResponse,
	}, nil
}
