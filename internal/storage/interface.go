package storage

import "context"

type Store interface {
	ReadEvent(ctx context.Context, hash string) (string, error)
	WriteEvent(ctx context.Context, hash, originalURL string) error
}
