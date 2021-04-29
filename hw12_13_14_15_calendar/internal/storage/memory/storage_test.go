package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage"

	"github.com/stretchr/testify/suite"
)

type MemStoreSuite struct {
	suite.Suite
	store *Storage
}

func (m *MemStoreSuite) SetupTest() {
	m.store = New()
	notification := 4 * time.Hour

	m.store.events = map[int]*storage.Event{
		1: &storage.Event{
			ID:           1,
			Title:        "Title1",
			Start:        time.Now().Add(2 * time.Hour),
			Stop:         time.Now().Add(3 * time.Hour),
			Description:  "the event",
			UserID:       1,
			Notification: &notification,
		},
		2: &storage.Event{
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
	newEvent := storage.Event{
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

//func (m *MemStoreSuite) TestInsertNewEventWithFail() {
//	newEvent := storage.Event{
//		ID:          "1",
//		Title:       "Title6",
//		StartDate:   100500,
//		EndDate:     200200,
//		Description: "Description6",
//		OwnerID:     "",
//		RemindIn:    150,
//	}
//
//	err := m.store.NewEvent(context.Background(), newEvent)
//
//	m.Require().Error(err)
//	m.Require().EqualError(ErrEventAlreadyExist, err.Error())
//}
//
//func (m *MemStoreSuite) TestUpdateEventSuccess() {
//	toUpdate := storage.Event{
//		ID:          "1",
//		Title:       "TitleUpdated",
//		StartDate:   1,
//		EndDate:     2,
//		Description: "DescriptionUpdate",
//		OwnerID:     "",
//		RemindIn:    5,
//	}
//	err := m.store.UpdateEvent(context.Background(), toUpdate)
//
//	m.Require().NoError(err)
//
//	updated := m.store.events["1"]
//
//	m.Require().NotNil(updated)
//	m.Require().Equal(toUpdate.Title, updated.Title)
//	m.Require().Equal(toUpdate.Description, updated.Description)
//}
//
//func (m *MemStoreSuite) TestUpdateEventWithError() {
//	toUpdate := storage.Event{
//		ID:          "6",
//		Title:       "Title6",
//		StartDate:   100500,
//		EndDate:     200200,
//		Description: "Description6",
//		OwnerID:     "",
//		RemindIn:    150,
//	}
//
//	err := m.store.UpdateEvent(context.Background(), toUpdate)
//
//	m.Require().Error(err)
//	m.Require().EqualError(ErrEventDoesNotExist, err.Error())
//}
//
//func (m *MemStoreSuite) TestRemoveEventSuccess() {
//	err := m.store.RemoveEvent(context.Background(), "1")
//
//	m.Require().NoError(err)
//
//	deleted := m.store.events["1"]
//
//	m.Require().Nil(deleted)
//}
//
//func (m *MemStoreSuite) TestRemoveEventWithError() {
//	err := m.store.RemoveEvent(context.Background(), "NaN")
//
//	m.Require().Error(err)
//	m.Require().EqualError(ErrEventDoesNotExist, err.Error())
//}
//
//func (m *MemStoreSuite) TestEventListSuccess() {
//	list, err := m.store.EventList(context.Background(), 3, 10)
//
//	m.Require().NoError(err)
//	m.Require().Len(list, 3)
//	m.Require().Contains(list, *m.store.events["2"])
//	m.Require().Contains(list, *m.store.events["3"])
//	m.Require().Contains(list, *m.store.events["4"])
//}
//
//func (m *MemStoreSuite) TestEventListWithError() {
//	list, err := m.store.EventList(context.Background(), 500, 700)
//
//	m.Require().Error(err)
//	m.Require().EqualError(ErrNoEvents, err.Error())
//	m.Require().Nil(list)
//}
//
//func (m *MemStoreSuite) TestAsyncOperations() {
//	var wg sync.WaitGroup
//
//	wg.Add(4)
//	go func() {
//		defer wg.Done()
//		for i := 10; i < 20; i++ {
//			newEvent := storage.Event{
//				ID:    fmt.Sprint(i),
//				Title: fmt.Sprintf("Title%d", i),
//			}
//
//			err := m.store.NewEvent(context.Background(), newEvent)
//
//			m.Require().NoError(err)
//
//			time.Sleep(30 * time.Millisecond)
//		}
//	}()
//
//	go func() {
//		defer wg.Done()
//		time.Sleep(200 * time.Millisecond)
//		toUpdate := storage.Event{
//			ID:          "1",
//			Title:       "TitleUpdate",
//			StartDate:   1,
//			EndDate:     2,
//			Description: "DescriptionUpdate",
//			OwnerID:     "",
//			RemindIn:    5,
//		}
//
//		err := m.store.UpdateEvent(context.Background(), toUpdate)
//
//		m.Require().NoError(err)
//	}()
//
//	go func() {
//		defer wg.Done()
//		time.Sleep(150 * time.Millisecond)
//		list, err := m.store.EventList(context.Background(), 6, 10)
//
//		m.Require().NoError(err)
//		m.Require().Len(list, 2)
//	}()
//
//	go func() {
//		defer wg.Done()
//		err := m.store.RemoveEvent(context.Background(), "5")
//		m.Require().NoError(err)
//	}()
//
//	wg.Wait()
//}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(MemStoreSuite))
}
