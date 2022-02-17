package integrationt_test

import (
	"context"
	"fmt"
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
	secondURL = "/resize/100/100/"
	imageURL  = "http://sebweo.com/wp-content/uploads/2020/01/what-is-jpeg_thumb-800x478.jpg?x72922"
)

var (
	Port string
	Host string
)

func setUp() {
	logg := logger.New()
	cfg := config.New()
	Host = cfg.App.Host
	Port = cfg.App.Port
	c := cache.New(cfg.Cache.Length)
	service := services.New(logg)
	imagePreviewer := app.New(logg, c, service)

	server := internalhttp.NewServer(logg, imagePreviewer, cfg)
	server.Start(context.Background())
}

func TestServer(t *testing.T) {
	serverReady := make(chan struct{})
	go func() {
		setUp()
		serverReady <- struct{}{}
	}()
	<-serverReady

	// second request with 100x100
	req, err := http.NewRequest(echo.GET, "http://"+Host+":"+Port+firstURL+imageURL, nil)

	require.NoError(t, err)
	client := http.Client{}
	fmt.Println(1)
	resp, err := client.Do(req)
	fmt.Println(2)

	require.NoError(t, err)

	require.Equal(t, resp.Header.Get("Content-Type"), "image/jpeg")

	firstByteBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	err = resp.Body.Close()
	require.NoError(t, err)

	// second request with 200x100
	req, err = http.NewRequest(echo.GET, "http://"+Host+":"+Port+secondURL+imageURL, nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, resp.Header.Get("Content-Type"), "image/jpeg")

	secondByteBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Greater(t, len(secondByteBody), len(firstByteBody))
	err = resp.Body.Close()
	require.NoError(t, err)
}
