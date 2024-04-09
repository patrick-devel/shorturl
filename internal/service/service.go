package service

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/patrick-devel/shorturl/internal/models"
)

type FileManager struct {
	consumer Consumer
	producer Producer
}

func New(consumer Consumer, producer Producer) *FileManager {
	return &FileManager{consumer: consumer, producer: producer}
}

type Consumer interface {
	ReadEvent(hash string) (*models.Event, error)
	Close() error
}
type Producer interface {
	WriteEvent(event *models.Event) error
	Close() error
}

func (fm *FileManager) ReadEvent(hash string) (string, error) {
	event, err := fm.consumer.ReadEvent(hash)
	if err != nil {
		return "", fmt.Errorf("error read event: %w", err)
	}

	return event.OriginalUrl, nil
}

func (fm *FileManager) WriteEvent(hash, originalUrl string) error {
	event := models.Event{UUID: uuid.NewString(), ShortUrl: hash, OriginalUrl: originalUrl}
	err := fm.producer.WriteEvent(&event)
	if err != nil {
		return fmt.Errorf("error write event: %w", err)
	}

	return nil
}
