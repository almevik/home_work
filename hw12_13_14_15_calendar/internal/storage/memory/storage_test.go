package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage/repository"
	"github.com/stretchr/testify/suite"
)

type MemStoreSuite struct {
	suite.Suite
	store *Storage
}

func (m *MemStoreSuite) SetupTest() {
	m.store = New()
	notification := 4 * time.Hour

	m.store.events = map[int]*repository.Event{
		1: &repository.Event{
			ID:           1,
			Title:        "Title1",
			Start:        time.Now().Add(2 * time.Hour),
			Stop:         time.Now().Add(3 * time.Hour),
			Description:  "the event",
			UserID:       1,
			Notification: &notification,
		},
		2: &repository.Event{
			ID:           2,
			Title:        "Title2",
			Start:        time.Now().Add(time.Hour),
			Stop:         time.Now().Add(2 * time.Hour),
			Description:  "Description2",
			UserID:       2,
			Notification: &notification,
		},
	}
}

func (m *MemStoreSuite) TestInsertNewEventSuccess() {
	newEvent := repository.Event{
		Title:       "Title6",
		Start:       time.Now().Add(time.Hour),
		Stop:        time.Now().Add(2 * time.Hour),
		Description: "Description6",
		UserID:      2,
	}

	id, err := m.store.CreateEvent(context.Background(), newEvent)

	m.Require().NoError(err)

	saved := m.store.events[1]

	m.Require().NotNil(saved)
	m.Require().Equal(1, id)
	m.Require().Equal(newEvent.Title, saved.Title)
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(MemStoreSuite))
}
