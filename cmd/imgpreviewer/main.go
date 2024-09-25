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
	"github.com/yakuninmax/imgpreviewer/internal/server"
	"github.com/yakuninmax/imgpreviewer/internal/storage"
)

func main() {
	logg, err := logger.New()
	if err != nil {
		log.Fatalln(err.Error())
	}

	conf, err := config.New(logg)
	if err != nil {
		logg.Error(err.Error())
		os.Exit(1)
	}

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

	cache := cache.New(conf.CacheSize(), store)

	dl := downloader.New(conf.RequestTimeout())

	app := app.New(logg, cache, dl)

	srv := server.New(conf.Port(), app, logg)

	go func() {
		logg.Info("starting server")
		err := srv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Error(err.Error())
			os.Exit(1)
		}

		logg.Info("stopped serving new connections")
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	sctx, rctx := context.WithTimeout(context.Background(), 5*time.Second)
	defer rctx()

	if err := srv.Stop(sctx); err != nil {
		logg.Error(err.Error())
		panic(err.Error())
	}
	logg.Info("graceful shutdown complete")
}
