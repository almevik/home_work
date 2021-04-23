package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
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
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go listenSignals(cancel)

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("failed read config %v\n", err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.FilePath)
	if err != nil {
		log.Fatalf("failed start logger %v\n", err)
	}

	logg.Info("init store...")

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

	logg.Info("init store completed...")
	logg.Info("starting calendar...")

	calendar := app.New(logg, store)

	server := internalhttp.NewServer(*calendar, logg, config.Server.Host, config.Server.Port)

	logg.Info("run server...")
	go func() {
		if err := server.Start(); err != nil {
			logg.Error(err)
			cancel()
		}
	}()
	logg.Info("server is running")

	<-ctx.Done()

	logg.Info("stopping server...")
	cancel()

	ctxSrv, cancelSrv := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelSrv()

	if err := server.Stop(ctxSrv); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}

	logg.Info("server is stopped")
}

func listenSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	cancel()
}
