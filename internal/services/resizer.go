package services

import (
	"bytes"
	"errors"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Georgiiagon/image-previewer/internal/app"
	"github.com/nfnt/resize"
)

var (
	ErrInvalidURL   = errors.New("invalid url")
	ErrFileNotImage = errors.New("invalid file type")
)

type Service struct {
	Logger app.Logger
}

func New(logger app.Logger) *Service {
	return &Service{
		Logger: logger,
	}
}

func (s *Service) Resize(byteImg []byte, height int, width int) ([]byte, string, error) {
	newName := s.RandString(30) + "-" + strconv.Itoa(height) + "-" + strconv.Itoa(width) + ".jpg"
	img, err := jpeg.Decode(bytes.NewBuffer(byteImg))
	if err != nil {
		return nil, "", err
	}

	img = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	out, err := os.Create("/tmp/" + newName)
	if err != nil {
		return nil, "", err
	}
	defer out.Close()

	// write new image to file
	err = jpeg.Encode(out, img, nil)

	if err != nil {
		return nil, "", err
	}

	resizedData, err := os.ReadFile("/tmp/" + newName)
	if err != nil {
		return nil, "", err
	}

	return resizedData, newName, nil
}

func (s *Service) Proxy(url string, headers http.Header) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	for k, vv := range r.Header {
		for _, v := range vv {
			headers.Add(k, v)
		}
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return nil, ErrInvalidURL
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(data)
	s.Logger.Info("contentType = " + contentType)
	if contentType == "image/jpeg" {
		return data, nil
	}

	return nil, ErrFileNotImage
}
