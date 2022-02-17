package internalhttp

import (
	"net/http"
	"os"
	"strconv"

	"github.com/Georgiiagon/image-previewer/internal/app"
	"github.com/labstack/echo/v4"
)

type APIHandler struct {
	App    Application
	Logger Logger
}

func (api *APIHandler) Resize(c echo.Context) error {
	height, err := strconv.Atoi(c.Param("height"))
	if err != nil {
		api.Logger.Error(err.Error())
		return nil
	}

	width, err := strconv.Atoi(c.Param("width"))
	if err != nil {
		api.Logger.Error(err.Error())
		return nil
	}

	if height < 1 || width < 1 {
		api.Logger.Error("width or height are too small")
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	url := c.Param("url")

	path, ok := api.App.Get(app.Key(url + "-" + strconv.Itoa(height) + "-" + strconv.Itoa(width)))

	if ok {
		resizedData, err := os.ReadFile(path.(string))
		if err != nil {
			return err
		}
		api.Logger.Info("Image found")

		c.Response().WriteHeader(http.StatusOK)
		c.Response().Header().Set("Content-Type", "image/jpeg")
		c.Response().Write(resizedData) //nolint:errcheck

		return nil
	}

	byteImg, err := api.App.Proxy(c.Param("url"), c.Response().Header())
	if err != nil {
		api.Logger.Error(err.Error())
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	imgBytes, newName, err := api.App.Resize(byteImg, height, width)
	if err != nil {
		api.Logger.Error(err.Error())
		return nil
	}

	api.App.Set(app.Key(url+"-"+strconv.Itoa(height)+"-"+strconv.Itoa(width)), "/tmp/"+newName)

	api.Logger.Info("Image saved")

	c.Response().WriteHeader(http.StatusOK)
	c.Response().Header().Set("Content-Type", "image/jpeg")
	c.Response().Write(imgBytes) //nolint:errcheck

	return nil
}

func NewHandler(app Application, logger Logger) APIHandler {
	return APIHandler{
		App:    app,
		Logger: logger,
	}
}
