package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	sync.RWMutex
	events map[int]*storage.Event
	last int
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(_ context.Context, _ string) error {
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func (s *Storage) CreateEvent(_ context.Context, event storage.Event) (int, error) {
	s.Lock()
	defer s.Unlock()

	id := s.next()
	event.ID = id
	s.events[id] = &event

	return id, nil
}

func (s *Storage) UpdateEvent(_ context.Context, id int, change storage.Event) error {
	s.Lock()
	defer s.Unlock()

	event, ok := s.events[id]
	if !ok {
		return storage.ErrEventNotFound
	}

	event.Title = change.Title
	event.Start = change.Start
	event.Stop = change.Stop
	event.Description = change.Description
	event.Notification = change.Notification
	s.events[id] = event

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id int) error {
	s.Lock()
	defer s.Unlock()

	delete(s.events, id)
	return nil
}

func (s *Storage) DeleteAllEvents(_ context.Context) error {
	s.Lock()
	defer s.Unlock()

	s.events = make(map[int]*storage.Event)
	return nil
}

func (s *Storage) ShowDayEvents(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.Lock()
	defer s.Unlock()

	var result []storage.Event
	year, month, day := date.Date()

	for _, event := range s.events {
		eYear, eMonth, eDay := event.Start.Date()
		if eYear == year && eMonth == month && eDay == day {
			result = append(result, *event)
		}
	}

	return order(result), nil
}

func (s *Storage) ShowWeekEvents(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.Lock()
	defer s.Unlock()

	var result []storage.Event
	year, week := date.ISOWeek()

	for _, event := range s.events {
		eYear, eWeek := event.Start.ISOWeek()
		if eYear == year && eWeek == week {
			result = append(result, *event)
		}
	}

	return order(result), nil
}

func (s *Storage) ShowMonthEvents(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.Lock()
	defer s.Unlock()

	var result []storage.Event
	year, month, _ := date.Date()

	for _, event := range s.events {
		eYear, eMonth, _ := event.Start.Date()
		if eYear == year && eMonth == month {
			result = append(result, *event)
		}
	}

	return order(result), nil
}

func (s *Storage) next() int {
	s.last++
	return s.last
}

// Сортировка событий по порядку.
func order(events []storage.Event) []storage.Event {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})

	return events
}
