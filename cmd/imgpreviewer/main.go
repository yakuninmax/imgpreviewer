package main

import (
	"log"
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

	headers := make(map[string]string)
	headers["Accept"] = "text/html"
	headers["Content-Type"] = "text/html; charset=utf-8"

	image, err := client.GetImage("https://github.com/notepad-plus-plus/notepad-plus-plus/releases/download/v8.6.7/npp.8.6.7.Installer.x64.exe", headers)
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}

	err = os.WriteFile("destination.jpeg", image, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
