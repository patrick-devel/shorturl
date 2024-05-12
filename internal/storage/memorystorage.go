package storage

import (
	"context"
	"fmt"

	"github.com/patrick-devel/shorturl/internal/models"
)

type MemoryStorage struct {
	cache map[string]string
}

func NewMemoryStorage(cache map[string]string) *MemoryStorage {
	return &MemoryStorage{cache: cache}
}

func (s *MemoryStorage) ReadEvent(_ context.Context, shortURL string) (string, error) {
	originalURL, ok := s.cache[shortURL]
	if !ok {
		return "", fmt.Errorf("error fetch event from memory")
	}

	return originalURL, nil
}

func (s *MemoryStorage) WriteEvent(_ context.Context, event models.Event) error {
	s.cache[event.ShortURL] = event.OriginalURL
	return nil
}

func (s *MemoryStorage) WriteEvents(_ context.Context, events []models.Event) error {
	for _, e := range events {
		s.cache[e.ShortURL] = e.OriginalURL
	}

	return nil
}

func (s *MemoryStorage) ReadEventsByCreatorID(_ context.Context, _ string) ([]models.Event, error) {
	return []models.Event{}, nil
}
