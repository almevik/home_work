package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrNoRows        = errors.New("no rows found")
)

type Storage interface {
	Connection
	Events
}

type Connection interface {
	Connect(ctx context.Context, dsn string) error
	Close(ctx context.Context) error
}

type Event struct {
	ID           int
	Title        string
	Start        time.Time
	Stop         time.Time
	Description  string
	UserID       int
	Notification *time.Duration
}

type Events interface {
	CreateEvent(ctx context.Context, event Event) (int, error)
	UpdateEvent(ctx context.Context, id int, event Event) error
	DeleteEvent(ctx context.Context, id int) error
	ShowDayEvents(ctx context.Context, date time.Time) ([]Event, error)
	ShowWeekEvents(ctx context.Context, date time.Time) ([]Event, error)
	ShowMonthEvents(ctx context.Context, date time.Time) ([]Event, error)
}
