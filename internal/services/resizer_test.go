package services

import (
	"net/http"
	"testing"

	"github.com/Georgiiagon/image-previewer/internal/logger"
	"github.com/stretchr/testify/require"
)

const URL = "https://sebweo.com/wp-content/uploads/2020/01/what-is-jpeg_thumb-800x478.jpg?x72922"

func TestProxy(t *testing.T) {
	logg := logger.New()
	service := New(logg)
	headers := http.Header{}
	imgBytes, err := service.Proxy(URL, headers)

	require.NoError(t, err)
	require.Greater(t, len(imgBytes), 0)
	require.Greater(t, len(headers), 1)
	require.Equal(t, "image/jpeg", headers.Get("Content-Type"))

	newImgBytes, _, err := service.Resize(imgBytes, 100, 100)
	require.NoError(t, err)
	require.Greater(t, len(newImgBytes), 0)
}
