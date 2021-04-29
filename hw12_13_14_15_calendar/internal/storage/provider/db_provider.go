package provider

import (
	"context"
	"errors"

	"github.com/almevik/home_work/hw12_13_14_15_calendar/config"
	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

var (
	ErrEmptyTitle     = errors.New("empty title")
	ErrIncorrectDates = errors.New("time start > time stop")
	ErrStopPassed     = errors.New("stop time already passed")
	ErrEmptyUserID    = errors.New("empty user id")
)

type Models struct {
	events storage.EventProvider
}

type DataProvider struct {
	db     storage.Storage
	cfg    config.StorageConf
	models *Models
}

func NewDataProvider(ctx context.Context, storageCfg config.StorageConf) (*DataProvider, error) {
	var db storage.Storage
	if storageCfg.Inmemory {
		db = memorystorage.New()
	} else {
		db = sqlstorage.New()
		err := db.Connect(ctx, storageCfg.DSN())
		if err != nil {
			return nil, err
		}
	}

	models := &Models{
		events: storage.EventProvider{
			DB: db,
		},
	}

	p := &DataProvider{
		db:     db,
		cfg:    storageCfg,
		models: models,
	}

	return p, nil
}

func (dp *DataProvider) Disconnect(ctx context.Context) error {
	return dp.db.Close(ctx)
}
