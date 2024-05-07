package storage

import (
	"context"

	"github.com/patrick-devel/shorturl/internal/models"
)

type Store interface {
	ReadEvent(ctx context.Context, hash string) (string, error)
	WriteEvent(ctx context.Context, event models.Event) error
	WriteEvents(_ context.Context, events []models.Event) error
}
