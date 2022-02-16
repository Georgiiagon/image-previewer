package app

import "net/http"

type Key string

type App struct {
	Logger  Logger
	Cache   Cache
	Service Service
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type Service interface {
	RandString(i int) string
	Resize(byteImg []byte, height int, width int) ([]byte, string, error)
	Proxy(url string, headers http.Header) ([]byte, error)
}

func New(logger Logger, cache Cache, service Service) *App {
	return &App{Logger: logger, Cache: cache, Service: service}
}

func (a *App) Set(key Key, value interface{}) bool {
	return a.Cache.Set(key, value)
}

func (a *App) Get(key Key) (interface{}, bool) {
	return a.Cache.Get(key)
}

func (a *App) Clear() {
	a.Cache.Clear()
}

func (a *App) Resize(byteImg []byte, height int, width int) ([]byte, string, error) {
	return a.Service.Resize(byteImg, height, width)
}

func (a *App) Proxy(url string, headers http.Header) ([]byte, error) {
	return a.Service.Proxy(url, headers)
}
