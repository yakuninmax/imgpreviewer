package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/yakuninmax/imgpreviewer/internal/config"
	"github.com/yakuninmax/imgpreviewer/internal/storage"
)

func main() {
	// Create logger.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Get config.
	config, err := config.New(logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Set up cache storage.
	storage, err := storage.New(config.CachePath(), config.CacheSize(), logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		err := storage.Clean()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	fmt.Print(storage.Path())
}
