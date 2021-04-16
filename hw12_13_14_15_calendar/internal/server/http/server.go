package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
)

type Server struct {
	app    app.App
	logger logger.Logger
	addr   string
	srv    *http.Server
	router *mux.Router
}

type Application interface {
}

func NewServer(app app.App, logg logger.Logger, host string, port string) *Server {
	return &Server{
		app:    app,
		logger: logg,
		addr:   net.JoinHostPort(host, port),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.SetupRoutes()

	s.srv = &http.Server{
		Addr:         s.addr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server closed: %w", err)
	}

	<-ctx.Done()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	return err
}
