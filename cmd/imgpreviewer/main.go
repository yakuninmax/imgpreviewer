package main

import (
	"fmt"
	"os"

	"github.com/yakuninmax/imgpreviewer/internal/app"
	"github.com/yakuninmax/imgpreviewer/internal/cache"
	"github.com/yakuninmax/imgpreviewer/internal/config"
	"github.com/yakuninmax/imgpreviewer/internal/downloader"
	"github.com/yakuninmax/imgpreviewer/internal/logger"
	"github.com/yakuninmax/imgpreviewer/internal/processor"
)

func main() {
	// Init logger.
	l := logger.New()

	// Get config.
	conf, err := config.New(l)
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	// Init cache.
	c, err := cache.New(conf.CachePath(), conf.CacheSize(), l)
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		err := c.Clean(l)
		if err != nil {
			l.Error(err.Error())
			os.Exit(1)
		}
	}()

	// Init downloader.
	d := downloader.New(conf.RequestTimeout())

	// Init processor.
	p := processor.New()

	// Init app.
	app := app.New(l, c, d, p)

	headers := make(map[string]string)
	headers["Accept"] = "text/html"
	headers["Content-Type"] = "text/html; charset=utf-8"

	_, err = app.Crop("500/300/https://filesamples.com/samples/image/jpeg/sample_1920%C3%971280.jpeg", headers)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
