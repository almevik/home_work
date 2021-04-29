package storage

import (
	"context"
	"errors"
	"time"

	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage/provider"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrNoRows        = errors.New("no rows found")
)

type Event struct {
	ID           int
	Title        string
	Start        time.Time
	Stop         time.Time
	Description  string
	UserID       int
	Notification *time.Duration
}

type EventProvider struct {
	DB Storage
}

func (ep *EventProvider) CreateEvent(ctx context.Context, event Event) (int, error) {
	if event.Title == "" {
		return 0, provider.ErrEmptyTitle
	}
	if event.Start.After(event.Stop) {
		return 0, provider.ErrIncorrectDates
	}
	// Не даем создавать прошедшие события
	if time.Now().After(event.Stop) {
		return 0, provider.ErrStopPassed
	}
	if event.UserID == 0 {
		return 0, provider.ErrEmptyUserID
	}

	return ep.DB.CreateEvent(ctx, event)
}

func (ep *EventProvider) UpdateEvent(ctx context.Context, id int, event Event) error {
	if event.Title == "" {
		return provider.ErrEmptyTitle
	}
	if event.Start.After(event.Stop) {
		return provider.ErrIncorrectDates
	}
	// Не даем менять прошедшие события
	if time.Now().After(event.Stop) {
		return provider.ErrStopPassed
	}

	return ep.DB.UpdateEvent(ctx, id, event)
}

func (ep *EventProvider) DeleteEvent(ctx context.Context, id int) error {
	return ep.DB.DeleteEvent(ctx, id)
}

func (ep *EventProvider) ShowDayEvents(ctx context.Context, date time.Time) ([]Event, error) {
	return ep.DB.ShowDayEvents(ctx, date)
}

func (ep *EventProvider) ShowWeekEvents(ctx context.Context, date time.Time) ([]Event, error) {
	return ep.DB.ShowWeekEvents(ctx, date)
}

func (ep *EventProvider) ShowMonthEvents(ctx context.Context, date time.Time) ([]Event, error) {
	return ep.DB.ShowMonthEvents(ctx, date)
}
