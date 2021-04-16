package internalhttp

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	Routes() []Route
}

type Route struct {
	Name   string
	Method string
	Path   string
	Func   http.HandlerFunc
}

// SetupRoutes устанавливает роуты.
func (s *Server) SetupRoutes() {
	s.router = mux.NewRouter()
	s.router.HandleFunc("/", s.homeHandler).Methods("GET")
	s.router.HandleFunc("/hello-world", s.homeHandler).Methods("GET")
	http.Handle("/", s.router)
}

func (s *Server) homeHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("Hello, world\n"))
	if err != nil {
		s.logger.Error(fmt.Errorf("http write: %w", err))
	}
}
