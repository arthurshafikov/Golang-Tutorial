package api

import (
	"context"
	"log"
	"net"
	"time"

	grpcapi "github.com/thewolf27/hw12_13_14_15_calendar/internal/server/grpc/generated"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
)

type Application interface {
	CreateEvent(storage.Event) error
	UpdateEvent(storage.Event) error
	DeleteEvent(storage.Event) error
	ListEventsOnADay(time.Time) (storage.EventsSlice, error)
	ListEventsOnAWeek(time.Time) (storage.EventsSlice, error)
	ListEventsOnAMonth(time.Time) (storage.EventsSlice, error)
}

func RunGrpcServer(ctx context.Context, address string, app Application) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	s := CalendarService{
		app: app,
	}
	grpcserver := grpc.NewServer()

	grpcapi.RegisterCalendarServiceServer(grpcserver, &s)

	go func() {
		<-ctx.Done()

		grpcserver.Stop()
	}()

	if err := grpcserver.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
