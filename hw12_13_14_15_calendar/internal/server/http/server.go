package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Server *http.Server
	ctx    context.Context
	Logger Logger
	App    Application
	Host   string
	Port   string
}

type Logger interface { // TODO
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Application interface { // TODO
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

	finalHandler := http.HandlerFunc(handleHello)
	m.Handle("/hello", loggingMiddleware(finalHandler, s))

	if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.Server.Shutdown(ctx)
	return nil
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}
