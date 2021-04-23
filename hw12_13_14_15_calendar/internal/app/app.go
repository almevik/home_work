package app

import (
	"context"
	"errors"
	"time"

	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger
	Storage
}

type Methods interface {
	CreateEvent(ctx context.Context, event storage.Event) (id int, err error)
	UpdateEvent(ctx context.Context, id int, event storage.Event) error
	DeleteEvent(ctx context.Context, id int) error
	ShowDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	ShowWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	ShowMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
}

type Logger interface {
	logger.Logger
}

type Storage interface {
	storage.Storage
}

func New(logger logger.Logger, storage storage.Storage) *App {
	return &App{
		logger,
		storage,
	}
}

var (
	ErrEmptyTitle     = errors.New("empty title")
	ErrIncorrectDates = errors.New("time start > time stop")
	ErrStopPassed     = errors.New("stop time already passed")
	ErrEmptyUserID    = errors.New("empty user id")
)

func (a *App) CreateEvent(ctx context.Context, event storage.Event) (int, error) {
	if event.Title == "" {
		return 0, ErrEmptyTitle
	}
	if event.Start.After(event.Stop) {
		return 0, ErrIncorrectDates
	}
	// Не даем создавать прошедшие события
	if time.Now().After(event.Stop) {
		return 0, ErrStopPassed
	}
	if event.UserID == 0 {
		return 0, ErrEmptyUserID
	}

	return a.Storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, id int, event storage.Event) error {
	if event.Title == "" {
		return ErrEmptyTitle
	}
	if event.Start.After(event.Stop) {
		return ErrIncorrectDates
	}
	// Не даем менять прошедшие события
	if time.Now().After(event.Stop) {
		return ErrStopPassed
	}

	return a.Storage.UpdateEvent(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	return a.Storage.DeleteEvent(ctx, id)
}

func (a *App) ShowDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.Storage.ShowDayEvents(ctx, date)
}

func (a *App) ShowWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.Storage.ShowWeekEvents(ctx, date)
}

func (a *App) ShowMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.Storage.ShowMonthEvents(ctx, date)
}
