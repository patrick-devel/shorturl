package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/patrick-devel/shorturl/internal/models"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{db: db}
}

func (s *DBStorage) ReadEvent(ctx context.Context, hash string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE hash=$1;", hash)

	var OriginalURL string

	err := row.Scan(&OriginalURL)
	if err != nil {
		return "", fmt.Errorf("error fetch event from db: %w", err)
	}

	return OriginalURL, nil
}

func (s *DBStorage) WriteEvent(ctx context.Context, hash, originalURL string) error {
	event := models.Event{UUID: uuid.NewString(), ShortURL: hash, OriginalURL: originalURL}

	sqlStatement := `INSERT INTO urls (uuid, hash, original_url) VALUES ($1, $2, $3);`
	_, err := s.db.ExecContext(ctx, sqlStatement, event.UUID, event.ShortURL, event.OriginalURL)
	if err != nil {
		return fmt.Errorf("error write event to db: %w", err)
	}

	return nil
}
