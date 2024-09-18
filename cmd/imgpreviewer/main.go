package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yakuninmax/imgpreviewer/internal/app"
	"github.com/yakuninmax/imgpreviewer/internal/cache"
	"github.com/yakuninmax/imgpreviewer/internal/config"
	"github.com/yakuninmax/imgpreviewer/internal/downloader"
	"github.com/yakuninmax/imgpreviewer/internal/logger"
	"github.com/yakuninmax/imgpreviewer/internal/processor"
	"github.com/yakuninmax/imgpreviewer/internal/server"
	"github.com/yakuninmax/imgpreviewer/internal/storage"
)

func main() {
	// Init logger.
	logg, err := logger.New()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Get config.
	conf, err := config.New(logg)
	if err != nil {
		logg.Error(err.Error())
		os.Exit(1)
	}

	// Init cache storage.
	store, err := storage.New(conf.CachePath())
	if err != nil {
		logg.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		err := store.Clean()
		if err != nil {
			logg.Error(err.Error())
			os.Exit(1)
		}
		logg.Info("cache cleared")
	}()
	logg.Info("temp cache folder is " + conf.CachePath())

	// Init cache.
	cache, err := cache.New(conf.CacheSize(), store)
	if err != nil {
		logg.Error(err.Error())
		panic(err.Error())
	}

	// Init downloader.
	dl := downloader.New(conf.RequestTimeout())

	// Init processor.
	proc := processor.New()

	// Init app.
	app := app.New(logg, cache, dl, proc)

	// Init server.
	srv := server.New(conf.Port(), app, logg)

	go func() {
		logg.Info("starting server")
		err := srv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Error(err.Error())
			panic(err.Error())
		}

		logg.Info("stopped serving new connections")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownRelease()

	if err := srv.Stop(shutdownCtx); err != nil {
		logg.Error(err.Error())
		panic(err.Error())
	}
	logg.Info("graceful shutdown complete")
}
