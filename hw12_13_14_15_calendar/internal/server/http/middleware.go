package internalhttp

import (
	"fmt"
	"net/http"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		s.logger.Info(
			fmt.Sprintf("%s %s %s %s %d %s",
				requestAddr(r),
				r.Method,
				r.RequestURI,
				r.Proto,
				rw.code,
				userAgent(r),
			))
	})
}
