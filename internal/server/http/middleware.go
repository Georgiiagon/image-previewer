package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info(fmt.Sprintf("%s [%s] %s %s %s %v %v %s",
			r.RemoteAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL,
			r.Proto,
			200,
			r.ContentLength,
			r.UserAgent()))

		h(w, r)
	}
}
