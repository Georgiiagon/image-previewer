package internalhttp

import (
	"net/http"

	"github.com/Georgiiagon/image-previewer/internal/config"
)

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func (apiHandler) Resize(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func NewHandler(config.Config) (http.Handler, error) {
	return apiHandler{}, nil
}
