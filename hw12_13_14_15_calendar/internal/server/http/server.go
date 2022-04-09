package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/arthurshafikov/hw12_13_14_15_calendar/internal/server"
	"github.com/arthurshafikov/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(storage.Event) (int64, error)
	UpdateEvent(storage.Event) (int64, error)
	DeleteEvent(storage.Event) error
	ListEventsOnADay(time.Time) (storage.EventsSlice, error)
	ListEventsOnAWeek(time.Time) (storage.EventsSlice, error)
	ListEventsOnAMonth(time.Time) (storage.EventsSlice, error)
}

type serverResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

type Server struct {
	Server *http.Server
	ctx    context.Context
	Logger Logger
	App    Application
	Host   string
	Port   string
}

func NewServer(logger Logger, app Application, host string, port string) *Server {
	return &Server{
		Logger: logger,
		App:    app,
		Host:   host,
		Port:   port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.ctx = ctx

	m := http.NewServeMux()
	server := http.Server{Addr: fmt.Sprintf("%v:%v", s.Host, s.Port), Handler: m}
	s.Server = &server

	m.Handle("/", s.loggingMiddleware(http.HandlerFunc(s.home)))
	m.Handle("/create", s.loggingMiddleware(http.HandlerFunc(s.create)))
	m.Handle("/update", s.loggingMiddleware(http.HandlerFunc(s.update)))
	m.Handle("/delete", s.loggingMiddleware(http.HandlerFunc(s.delete)))
	m.Handle("/list-a-day", s.loggingMiddleware(http.HandlerFunc(s.listEventsOnADay)))
	m.Handle("/list-a-week", s.loggingMiddleware(http.HandlerFunc(s.listEventsOnAWeek)))
	m.Handle("/list-a-month", s.loggingMiddleware(http.HandlerFunc(s.listEventsOnAMonth)))

	if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.Server.Shutdown(ctx)
	return nil
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, string("OK"))
}

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	event, err := server.DecodeJSONEvent(r.Body, true)
	if err != nil {
		s.writeJSONError(w, err)
		return
	}
	if _, err = s.App.CreateEvent(event); err != nil {
		s.writeJSONError(w, fmt.Errorf(server.ErrCantCreateEventFormat, err))
		return
	}
	s.writeJSONSuccess(w)
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	event, err := server.DecodeJSONEvent(r.Body, false)
	if err != nil {
		s.writeJSONError(w, err)
		return
	}
	if _, err = s.App.UpdateEvent(event); err != nil {
		s.writeJSONError(w, fmt.Errorf(server.ErrCantUpdateEventFormat, err))
		return
	}
	s.writeJSONSuccess(w)
}

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	event, err := server.DecodeJSONEvent(r.Body, false)
	if err != nil {
		s.writeJSONError(w, err)
		return
	}
	if err = s.App.DeleteEvent(event); err != nil {
		s.writeJSONError(w, fmt.Errorf(server.ErrCantDeleteEventFormat, err))
		return
	}
	s.writeJSONSuccess(w)
}

func (s *Server) listEventsOnADay(w http.ResponseWriter, r *http.Request) {
	s.listEvents(w, r, s.App.ListEventsOnADay)
}

func (s *Server) listEventsOnAWeek(w http.ResponseWriter, r *http.Request) {
	s.listEvents(w, r, s.App.ListEventsOnAWeek)
}

func (s *Server) listEventsOnAMonth(w http.ResponseWriter, r *http.Request) {
	s.listEvents(w, r, s.App.ListEventsOnAMonth)
}

func (s *Server) listEvents(w http.ResponseWriter, r *http.Request, listEventsF storage.ListEventsFunction) {
	date, err := server.DecodeJSONDate(r.Body)
	if err != nil {
		s.writeJSONError(w, err)
		return
	}
	events, err := listEventsF(date)
	if err != nil {
		s.writeJSONError(w, fmt.Errorf(server.ErrCantListEventsFormat, err))
		return
	}
	s.writeJSONResponse(w, serverResponse{
		Data: events,
	})
}

func (s *Server) writeJSONError(w http.ResponseWriter, err error) {
	s.Logger.Error(err.Error())
	s.writeJSONResponse(w, serverResponse{
		Error: err.Error(),
	})
}

func (s *Server) writeJSONSuccess(w http.ResponseWriter) {
	s.writeJSONResponse(w, serverResponse{
		Data: storage.RequestSuccessMessage,
	})
}

func (s *Server) writeJSONResponse(w http.ResponseWriter, response serverResponse) {
	w.Header().Set("Content-Type", "application/json")
	status := http.StatusOK

	m, err := json.Marshal(response)
	if err != nil {
		m = []byte("Fatal error, can't marshal response")
	}

	if response.Data == nil || err != nil {
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)

	fmt.Fprint(w, string(m))
}
