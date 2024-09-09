package main

import (
	"os"

	"github.com/yakuninmax/imgpreviewer/internal/cache"
	"github.com/yakuninmax/imgpreviewer/internal/config"
	httpclient "github.com/yakuninmax/imgpreviewer/internal/http/client"
	"github.com/yakuninmax/imgpreviewer/internal/logger"
)

func main() {
	// Init logger.
	logger := logger.New()

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

	// Init http client.
	client := httpclient.New(config.RequestTimeout())

	client.GetImage("https://ya.ru")
}
