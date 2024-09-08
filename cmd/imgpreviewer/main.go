package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/yakuninmax/imgpreviewer/internal/cache"
	"github.com/yakuninmax/imgpreviewer/internal/config"
)

func main() {
	// Init logger.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Get config.
	config, err := config.New(logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Init cache.
	cache, err := cache.New(config.CachePath(), config.CacheSize(), logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		err := cache.Clean(logger)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	fmt.Print(cache)
}
