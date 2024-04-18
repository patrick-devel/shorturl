package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	filemanager "github.com/patrick-devel/shorturl/internal/file_manager"
	"github.com/patrick-devel/shorturl/internal/models"
)

type FileStorage struct {
	consumer Consumer
	producer Producer
}

func NewFileStorage(path string) (*FileStorage, error) {
	consumer, err := filemanager.NewConsumer(path)
	if err != nil {
		return nil, err
	}

	producer, err := filemanager.NewProducer(path)
	if err != nil {
		return nil, err
	}
	return &FileStorage{consumer: consumer, producer: producer}, nil
}

type Consumer interface {
	ReadEvent(hash string) (*models.Event, error)
	Close() error
}

type Producer interface {
	WriteEvent(event *models.Event) error
	Close() error
}

func (fs *FileStorage) ReadEvent(_ context.Context, hash string) (string, error) {
	event, err := fs.consumer.ReadEvent(hash)
	if err != nil {
		return "", fmt.Errorf("error read event: %w", err)
	}

	return event.OriginalURL, nil
}

func (fs *FileStorage) WriteEvent(_ context.Context, hash, originalURL string) error {
	event := models.Event{UUID: uuid.NewString(), ShortURL: hash, OriginalURL: originalURL}
	err := fs.producer.WriteEvent(&event)
	if err != nil {
		return fmt.Errorf("error write event: %w", err)
	}

	return nil
}
