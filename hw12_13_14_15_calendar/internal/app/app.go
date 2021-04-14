package app

import (
	"context"
	"errors"
	"time"

	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage"
)

type app struct {
	Logger
	Storage
}

type App interface {
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

func New(logger logger.Logger, storage storage.Storage) *app {
	return &app{
		logger,
		storage,
	}
}

var ErrEmptyTitle = errors.New("empty title")
var ErrIncorrectDates = errors.New("time start > time stop")
var ErrStopPassed = errors.New("stop time already passed")
var ErrEmptyUserID = errors.New("empty user id")

func (a *app) CreateEvent(ctx context.Context, event storage.Event) (int, error) {
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

func (a *app) UpdateEvent(ctx context.Context, id int, event storage.Event) error {
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

func (a *app) DeleteEvent(ctx context.Context, id int) error {
	return a.Storage.DeleteEvent(ctx, id)
}

func (a *app) ShowDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.Storage.ShowDayEvents(ctx, date)
}

func (a *app) ShowWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.Storage.ShowWeekEvents(ctx, date)
}

func (a *app) ShowMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.Storage.ShowMonthEvents(ctx, date)
}
