package repository

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEmptyTitle     = errors.New("empty title")
	ErrIncorrectDates = errors.New("time start > time stop")
	ErrStopPassed     = errors.New("stop time already passed")
	ErrEmptyUserID    = errors.New("empty user id")
)

type Storage interface {
	Connect(ctx context.Context, dsn string) error
	Close(ctx context.Context) error

	Events
}

type Events interface {
	CreateEvent(ctx context.Context, event Event) (int, error)
	UpdateEvent(ctx context.Context, id int, event Event) error
	DeleteEvent(ctx context.Context, id int) error
	ShowDayEvents(ctx context.Context, date time.Time) ([]Event, error)
	ShowWeekEvents(ctx context.Context, date time.Time) ([]Event, error)
	ShowMonthEvents(ctx context.Context, date time.Time) ([]Event, error)
}
