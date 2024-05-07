package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"

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

func (s *DBStorage) WriteEvent(ctx context.Context, event models.Event) error {
	sqlStatement := `INSERT INTO urls (uuid, hash, original_url) VALUES ($1, $2, $3);`
	_, err := s.db.ExecContext(ctx, sqlStatement, event.UUID, event.Hash, event.OriginalURL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrDuplicateURL
		}

		return fmt.Errorf("error write event to db: %w", err)
	}

	return nil
}

func (s *DBStorage) WriteEvents(ctx context.Context, events []models.Event) error {
	sqlStatement := `INSERT INTO urls (uuid, hash, original_url) VALUES ($1, $2, $3) ON CONFLICT (url) DO UPDATE SET uuid = EXCLUDED.uuid;`

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("tx error: %w", err)
	}

	for _, e := range events {
		_, err = tx.ExecContext(ctx, sqlStatement, e.UUID, e.Hash, e.OriginalURL)
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
