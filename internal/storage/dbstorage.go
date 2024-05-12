package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/patrick-devel/shorturl/internal/models"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{db: db}
}

func (s *DBStorage) ReadEvent(ctx context.Context, shortURL string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE short_url=$1;", shortURL)

	var OriginalURL string

	err := row.Scan(&OriginalURL)
	if err != nil {
		return "", fmt.Errorf("error fetch event from db: %w", err)
	}

	return OriginalURL, nil
}

func (s *DBStorage) WriteEvent(ctx context.Context, event models.Event) error {
	sqlStatement := `INSERT INTO urls (uuid, creator_id, short_url, original_url) VALUES ($1, $2, $3, $4);`
	_, err := s.db.ExecContext(ctx, sqlStatement, event.UUID, event.CreatorID, event.ShortURL, event.OriginalURL)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = ErrDuplicateURL
		}

		return fmt.Errorf("error write event to db: %w", err)
	}

	return nil
}

func (s *DBStorage) WriteEvents(ctx context.Context, events []models.Event) error {
	sqlStatement := `INSERT INTO urls (uuid, creator_id, short_url, original_url) VALUES ($1, $2, $3, $4) ON CONFLICT (original_url) DO UPDATE SET uuid = EXCLUDED.uuid;`

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("tx error: %w", err)
	}

	for _, e := range events {
		_, err = tx.ExecContext(ctx, sqlStatement, e.UUID, e.CreatorID, e.ShortURL, e.OriginalURL)
		if err != nil {
			if rbError := tx.Rollback(); rbError != nil {
				logrus.Errorf("insert failed, unable to rollback %v", rbError)
			}

			return fmt.Errorf("error write event to db: %w", err)
		}
	}

	if cError := tx.Commit(); cError != nil {
		return fmt.Errorf("commit error: %w", cError)
	}

	return nil
}

func (s *DBStorage) ReadEventsByCreatorID(ctx context.Context, userID string) ([]models.Event, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT uuid, creator_id, short_url, original_url FROM urls WHERE creator_id=$1;", userID)
	if err != nil {
		return []models.Event{}, fmt.Errorf("error fetch events from db: %w", err)
	}

	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		event := new(models.Event)
		if err := rows.Scan(&event.UUID, &event.CreatorID, &event.ShortURL, &event.OriginalURL); err != nil {
			return events, fmt.Errorf("error decode events from db: %w", err)
		}
		events = append(events, *event)

	}

	if rows.Err() != nil {
		return events, fmt.Errorf("error scan rows: %w", rows.Err())
	}

	return events, nil
}
