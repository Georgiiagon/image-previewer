package integrationt_test

import (
	"context"
	"io/ioutil"
	"net/http"
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
	firstURL  = "/resize/100/100/"
	secondURL = "/resize/200/100/"
	wrongURL  = "/resize/-1/100/"
	imageURL  = "http://sebweo.com/wp-content/uploads/2020/01/what-is-jpeg_thumb-800x478.jpg?x72922"
)

var (
	Port string
	Host string
)

func setUp(ch chan struct{}) {
	logg := logger.New()
	cfg := config.New()
	Host = cfg.App.Host
	Port = cfg.App.Port
	c := cache.New(cfg.Cache.Length)
	service := services.New(logg)
	imagePreviewer := app.New(logg, c, service)

	server := internalhttp.NewServer(logg, imagePreviewer, cfg)
	ch <- struct{}{}
	server.Start(context.Background())
}

func TestServer(t *testing.T) {
	serverReady := make(chan struct{})
	go func() {
		setUp(serverReady)
	}()
	<-serverReady

	// second request with 100x100
	req, err := http.NewRequest(echo.GET, "http://"+Host+":"+Port+firstURL+imageURL, nil)

	require.NoError(t, err)
	client := http.Client{}
	resp, err := client.Do(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, resp.Header.Get("Content-Type"), "image/jpeg")

	firstByteBody, _ := ioutil.ReadAll(resp.Body)

	err = resp.Body.Close()
	require.NoError(t, err)

	// second request with 200x100
	req, err = http.NewRequest(echo.GET, "http://"+Host+":"+Port+secondURL+imageURL, nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, resp.Header.Get("Content-Type"), "image/jpeg")

	secondByteBody, _ := ioutil.ReadAll(resp.Body)

	require.Greater(t, len(secondByteBody), len(firstByteBody))
	err = resp.Body.Close()
	require.NoError(t, err)
}

func TestServerError(t *testing.T) {
	serverReady := make(chan struct{})
	go func() {
		setUp(serverReady)
	}()
	<-serverReady

	// second request with 100x100
	req, err := http.NewRequest(echo.GET, "http://"+Host+":"+Port+wrongURL+imageURL, nil)

	require.NoError(t, err)
	client := http.Client{}
	resp, err := client.Do(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.Equal(t, resp.Header.Get("Content-Type"), "")
	err = resp.Body.Close()
	require.NoError(t, err)
}

func TestWrongUrl(t *testing.T) {
	serverReady := make(chan struct{})
	go func() {
		setUp(serverReady)
	}()
	<-serverReady

	// second request with 100x100
	req, err := http.NewRequest(echo.GET, "http://"+Host+":"+Port+firstURL+"someurl", nil)

	require.NoError(t, err)
	client := http.Client{}
	resp, err := client.Do(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.Equal(t, resp.Header.Get("Content-Type"), "")
	err = resp.Body.Close()
	require.NoError(t, err)
}
