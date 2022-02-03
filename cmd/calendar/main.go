package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Georgiiagon/image-previewer/internal/app"
	"github.com/Georgiiagon/image-previewer/internal/config"
	"github.com/Georgiiagon/image-previewer/internal/logger"
	internalhttp "github.com/Georgiiagon/image-previewer/internal/server/http"
	"github.com/Georgiiagon/image-previewer/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", ".configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	fmt.Println(configFile)
	config := config.NewConfig(configFile)
	logg := logger.New(config.Logger.Level)

	storage := storage.New(config.Database)
	imagePreviewer := app.New(logg, storage)

	server := internalhttp.NewServer(logg, imagePreviewer, config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
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
