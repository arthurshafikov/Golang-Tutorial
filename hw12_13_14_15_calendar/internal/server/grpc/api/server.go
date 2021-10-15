package api

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/server/grpc/generated"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
)

type Application interface {
	CreateEvent(storage.Event) (int64, error)
	UpdateEvent(storage.Event) (int64, error)
	DeleteEvent(storage.Event) error
	ListEventsOnADay(time.Time) (storage.EventsSlice, error)
	ListEventsOnAWeek(time.Time) (storage.EventsSlice, error)
	ListEventsOnAMonth(time.Time) (storage.EventsSlice, error)
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func RunGrpcServer(ctx context.Context, host string, port string, app Application, logger Logger) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	s := CalendarService{
		app:    app,
		logger: logger,
	}
	grpcserver := grpc.NewServer(s.withServerUnaryInterceptor())
	generated.RegisterCalendarServiceServer(grpcserver, &s)

	go func() {
		<-ctx.Done()
		grpcserver.Stop()
	}()

	if err := grpcserver.Serve(lis); err != nil {
		logger.Error(err.Error())
		return
	}
}

func (c *CalendarService) serverInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	h, err := handler(ctx, req)

	c.logger.Info(
		fmt.Sprintf("GRPC Query [%v] %v",
			time.Now().Format("02/01/2006:15:04:05 MST"),
			info.FullMethod,
		),
	)

	return h, err
}

func (c *CalendarService) withServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(c.serverInterceptor)
}
