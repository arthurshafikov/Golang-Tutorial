package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info(
			fmt.Sprintf("%v [%v] %v %v %v",
				r.RemoteAddr,
				time.Now().Format("02/01/2006:15:04:05 MST"),
				r.Method,
				r.RequestURI,
				r.UserAgent(),
			),
		)
		next.ServeHTTP(w, r)
	})
}
