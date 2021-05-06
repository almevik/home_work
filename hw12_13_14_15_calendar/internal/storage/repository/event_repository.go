package repository

import (
	"context"
	"errors"
	"time"
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
	Ctx context.Context
	DB  Storage
}

func (ep *EventProvider) CreateEvent(event Event) (int, error) {
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

	return ep.DB.CreateEvent(ep.Ctx, event)
}

func (ep *EventProvider) UpdateEvent(id int, event Event) error {
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

	return ep.DB.UpdateEvent(ep.Ctx, id, event)
}

func (ep *EventProvider) DeleteEvent(id int) error {
	return ep.DB.DeleteEvent(ep.Ctx, id)
}

func (ep *EventProvider) ShowDayEvents(date time.Time) ([]Event, error) {
	return ep.DB.ShowDayEvents(ep.Ctx, date)
}

func (ep *EventProvider) ShowWeekEvents(date time.Time) ([]Event, error) {
	return ep.DB.ShowWeekEvents(ep.Ctx, date)
}

func (ep *EventProvider) ShowMonthEvents(date time.Time) ([]Event, error) {
	return ep.DB.ShowMonthEvents(ep.Ctx, date)
}
