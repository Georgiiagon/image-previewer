package internalhttp

import (
	"context"
	"net/http"

	"github.com/Georgiiagon/image-previewer/internal/app"
	"github.com/Georgiiagon/image-previewer/internal/config"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Logger Logger
	Config config.Config
	App    Application
	Server *echo.Echo
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface {
	Set(key app.Key, value interface{}) bool
	Get(key app.Key) (interface{}, bool)
	Clear()
	Resize(byteImg []byte, length int, width int) ([]byte, string, error)
	Proxy(url string, headers http.Header) ([]byte, error)
}

func NewServer(logger Logger, app Application, cfg config.Config) *Server {
	return &Server{
		Logger: logger,
		Config: cfg,
		App:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	handler := NewHandler(s.App, s.Logger)
	s.Server = echo.New()

	s.Server.GET("/resize/:height/:width/:url", handler.Resize)
	s.Server.Use(s.loggingMiddleware)

	ch := make(chan error)
	go func() {
		err := s.Server.Start(":" + s.Config.App.Port)
		ch <- err
	}()

	select {
	case <-ctx.Done():
		s.Logger.Info("Closed by context")
		err := s.Server.Shutdown(ctx)
		if err != nil {
			return err
		}
	case err := <-ch:
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.Logger.Info("Stop http server")

	return s.Server.Shutdown(ctx)
}
