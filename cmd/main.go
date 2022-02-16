package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Georgiiagon/image-previewer/internal/app"
	"github.com/Georgiiagon/image-previewer/internal/cache"
	"github.com/Georgiiagon/image-previewer/internal/config"
	"github.com/Georgiiagon/image-previewer/internal/logger"
	internalhttp "github.com/Georgiiagon/image-previewer/internal/server/http"
	"github.com/Georgiiagon/image-previewer/internal/services"
)

func main() {
	logg := logger.New()
	cfg := config.New()
	c := cache.New(cfg.Cache.Length)
	service := services.New(logg)
	imagePreviewer := app.New(logg, c, service)

	server := internalhttp.NewServer(logg, imagePreviewer, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		if err := server.Stop(context.Background()); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("image-previewer is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
