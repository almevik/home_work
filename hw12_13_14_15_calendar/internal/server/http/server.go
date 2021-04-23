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
	"strings"
	"time"
)

type Server struct {
	addr   string
	app    app.App
	logger logger.Logger
	srv    *http.Server
	router *mux.Router
}

type Application interface {
}

func NewServer(app app.App, logg logger.Logger, host string, port string) *Server {
	s := &Server{
		addr:   net.JoinHostPort(host, port),
		app:    app,
		logger: logg,
		router: mux.NewRouter(),
	}
	s.configureRouter()
	return s
}

func (s *Server) Start() error {
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

	return err
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	return err
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// handlers устанавливает роуты.
func (s *Server) configureRouter() {
	s.router.Use(s.loggingMiddleware)
	s.router.HandleFunc("/", s.homeHandler).Methods("GET")
	s.router.HandleFunc("/hello-world", s.homeHandler).Methods("GET")
}

func requestAddr(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func userAgent(r *http.Request) string {
	userAgents := r.Header["User-Agent"]
	if len(userAgents) > 0 {
		return "\"" + userAgents[0] + "\""
	}
	return ""
}

func (s *Server) homeHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("Hello, world\n"))
	if err != nil {
		s.logger.Error(fmt.Errorf("http write: %w", err))
	}
}
