package app

import (
	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage/provider"
)

type App struct {
	Logger
	Provider *provider.DataProvider
}

type Logger interface {
	logger.Logger
}

func New(logger logger.Logger, provider *provider.DataProvider) *App {
	return &App{
		logger,
		provider,
	}
}
