package storage

import (
	"context"
	"fmt"
)

type MemoryStorage struct {
	cache map[string]string
}

func NewMemoryStorage(cache map[string]string) *MemoryStorage {
	return &MemoryStorage{cache: cache}
}

func (s *MemoryStorage) ReadEvent(_ context.Context, hash string) (string, error) {
	originalURL, ok := s.cache[hash]
	if !ok {
		return "", fmt.Errorf("error fetch event from memory")
	}

	return originalURL, nil
}

func (s *MemoryStorage) WriteEvent(_ context.Context, hash, originalURL string) error {
	s.cache[hash] = originalURL
	return nil
}
