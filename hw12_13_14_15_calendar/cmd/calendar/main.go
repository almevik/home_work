package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/almevik/home_work/hw12_13_14_15_calendar/internal/server/http"
	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("failed read config", err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.FilePath)
	if err != nil {
		log.Fatalf("failed start logger %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var store storage.Storage

	if config.Storage.Inmemory {
		store = memorystorage.New()
	} else {
		dsn := DSN(config.Storage.Database)
		store = sqlstorage.New()

		err := store.Connect(ctx, dsn)
		if err != nil {
			logg.Error("failed connect storage" + err.Error())
		}
	}

	calendar := app.New(logg, store)

	server := internalhttp.NewServer(calendar, logg, config.Server.Host, config.Server.Port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
