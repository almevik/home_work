package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/almevik/home_work/hw12_13_14_15_calendar/internal/storage"
)

func New() *Storage {
	return &Storage{}
}

type Storage struct {
	db     *sql.DB
	logger logger.Logger
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect mysql: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute)

	s.db = db
	return s.db.PingContext(ctx)
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (int, error) {
	query := `INSERT INTO event (title, start, stop, description, user_id, notification)
			  VALUES(?, ?, ?, ?, ?, ?)`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			s.logger.Error(err)
		}
	}()

	res, err := stmt.Exec(
		event.Title,
		event.Start,
		event.Stop,
		event.Description,
		event.UserID,
		event.Notification,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id int, event storage.Event) error {
	query := `UPDATE event
				SET title = ?,
					start = ?,
					stop = ?,
					description = ?,
					notification = ?
				WHERE event_id = ?;`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			s.logger.Error(err)
		}
	}()

	res, err := stmt.Exec(
		event.Title,
		event.Start,
		event.Stop,
		event.Description,
		event.UserID,
		event.Notification,
		id,
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return storage.ErrEventNotFound
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
	query := `DELETE FROM event WHERE event_id = ?`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			s.logger.Error(err)
		}
	}()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return storage.ErrEventNotFound
	}
	return nil
}

func (s *Storage) DeleteAllEvents(ctx context.Context) error {
	query := `TRUNCATE event`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			s.logger.Error(err)
		}
	}()

	_, err = stmt.Exec()

	return err
}

func (s *Storage) ShowDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	y, m, d := date.Date()
	query := `SELECT id, title, start, stop, description, user_id, notification
		FROM event
		WHERE YEAR(start) = ? AND MONTH(start) = ? AND DAY(start) = ?
		ORDER BY start`

	var args []interface{}
	args = append(args, y, m, d)

	return s.searchEvents(ctx, query, args)
}

func (s *Storage) ShowWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	y, w := date.ISOWeek()
	query := `SELECT id, title, start, stop, description, user_id, notification
		FROM event
		WHERE YEAR(start) = ? AND WEEK(start) = ?
		ORDER BY start`

	var args []interface{}
	args = append(args, y, w)

	return s.searchEvents(ctx, query, args)
}

func (s *Storage) ShowMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	y, m, _ := date.Date()
	query := `SELECT id, title, start, stop, description, user_id, notification
		FROM event
		WHERE YEAR(start) = ? AND MONTH(start) = ?
		ORDER BY start`

	var args []interface{}
	args = append(args, y, m)

	return s.searchEvents(ctx, query, args)
}

// Общий запрос для поиска событий.
func (s *Storage) searchEvents(ctx context.Context, query string, args ...interface{}) ([]storage.Event, error) {
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			s.logger.Error(err)
		}
	}()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	var events []storage.Event
	for rows.Next() {
		event := new(storage.Event)
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Start,
			&event.Stop,
			&event.Description,
			&event.UserID,
			&event.Notification,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	if len(events) == 0 {
		return nil, storage.ErrNoRows
	}

	return events, nil
}
