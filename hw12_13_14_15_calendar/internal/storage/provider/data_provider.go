package provider

import (
	"context"

	"github.com/almevik/home_work/hw12_13_14_15_calendar/config"
	memorystorage "github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage/repository"
	sqlstorage "github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

type Models struct {
	Events repository.EventProvider
}

type DataProvider struct {
	db     repository.Storage
	cfg    config.StorageConf
	Models *Models
}

func NewDataProvider(ctx context.Context, storageCfg config.StorageConf) (*DataProvider, error) {
	var db repository.Storage
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
		Events: repository.EventProvider{
			Ctx: ctx,
			DB:  db,
		},
	}

	p := &DataProvider{
		db:     db,
		cfg:    storageCfg,
		Models: models,
	}

	return p, nil
}
