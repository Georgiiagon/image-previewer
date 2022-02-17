package internalhttp

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *Server) loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s.Logger.Info(fmt.Sprintf("%s [%s] %s %s %s %v %v %s",
			c.Request().URL.Path,
			c.Request().RemoteAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			c.Request().Method,
			c.Path(),
			c.Request().Proto,
			c.Request().ContentLength,
			c.Request().UserAgent()))
		return next(c)
	}
}
