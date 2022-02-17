package tests

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Georgiiagon/image-previewer/internal/app"
	"github.com/Georgiiagon/image-previewer/internal/cache"
	"github.com/Georgiiagon/image-previewer/internal/config"
	"github.com/Georgiiagon/image-previewer/internal/logger"
	internalhttp "github.com/Georgiiagon/image-previewer/internal/server/http"
	"github.com/Georgiiagon/image-previewer/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

const (
	ImageURL = "https://sebweo.com/wp-content/uploads/2020/01/what-is-jpeg_thumb-800x478.jpg?x72922"
	WrongURL = "https://sebweo.cqqm/umb-800x478.jpg?x72922"
)

func TestResize(t *testing.T) {
	logg := logger.New()
	cfg := config.New()
	c := cache.New(cfg.Cache.Length)
	service := services.New(logg)
	imagePreviewer := app.New(logg, c, service)

	// test lru cache
	_, ok := imagePreviewer.Get(app.Key(ImageURL + "-" + strconv.Itoa(100) + "-" + strconv.Itoa(200)))
	require.False(t, ok)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	con := e.NewContext(req, rec)
	con.SetPath("/resize/:height/:width/:url")
	con.SetParamNames("height", "width", "url")
	con.SetParamValues("100", "200", ImageURL)
	h := internalhttp.NewHandler(imagePreviewer, logg)
	require.NoError(t, h.Resize(con))
	require.Equal(t, http.StatusOK, rec.Code)

	require.Equal(t, rec.Header().Get("Content-Type"), "image/jpeg")
	contentLenght, err := strconv.Atoi(rec.Header().Get("Content-Length"))
	require.NoError(t, err)
	require.Greater(t, contentLenght, 50000)

	// test lru cache
	_, ok = imagePreviewer.Get(app.Key(ImageURL + "-" + strconv.Itoa(100) + "-" + strconv.Itoa(200)))
	require.True(t, ok)
}

func TestNegativeHighResize(t *testing.T) {
	logg := logger.New()
	cfg := config.New()
	c := cache.New(cfg.Cache.Length)
	service := services.New(logg)
	imagePreviewer := app.New(logg, c, service)

	// test lru cache
	_, ok := imagePreviewer.Get(app.Key(WrongURL + "-" + strconv.Itoa(100) + "-" + strconv.Itoa(200)))
	require.False(t, ok)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	con := e.NewContext(req, rec)
	con.SetPath("/resize/:height/:width/:url")
	con.SetParamNames("height", "width", "url")
	con.SetParamValues("-1", "200", WrongURL)

	h := internalhttp.NewHandler(imagePreviewer, logg)
	require.NoError(t, h.Resize(con))
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, rec.Header().Get("Content-Type"), "")
}

func TestServerNotFound(t *testing.T) {
	logg := logger.New()
	cfg := config.New()
	c := cache.New(cfg.Cache.Length)
	service := services.New(logg)
	imagePreviewer := app.New(logg, c, service)

	// test lru cache
	_, ok := imagePreviewer.Get(app.Key(WrongURL + "-" + strconv.Itoa(100) + "-" + strconv.Itoa(200)))
	require.False(t, ok)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	con := e.NewContext(req, rec)
	con.SetPath("/resize/:height/:width/:url")
	con.SetParamNames("height", "width", "url")
	con.SetParamValues("200", "200", WrongURL)

	h := internalhttp.NewHandler(imagePreviewer, logg)
	require.NoError(t, h.Resize(con))
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
