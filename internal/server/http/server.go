package internalhttp

import (
	"context"
	"net/http"

	"github.com/Georgiiagon/image-previewer/internal/config"
)

type Server struct {
	Logger Logger
	Config config.Config
	App    Application
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(context.Context, string, string) error
	GetEvents(ctx context.Context) error
}

func NewServer(logger Logger, app Application, conf config.Config) *Server {
	return &Server{
		Logger: logger,
		Config: conf,
		App:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	handler := apiHandler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/resize", s.loggingMiddleware(handler.Resize))

	ch := make(chan error)
	go func() {
		err := http.ListenAndServe(s.Config.App.Host+":"+s.Config.App.Port, mux)
		ch <- err
	}()

	select {
	case <-ctx.Done():
		s.Logger.Info("Closed by context")
	case err := <-ch:
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.Logger.Info("Stop http server")
	return nil
}
