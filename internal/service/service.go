package service

import (
	"context"
	"fmt"
	"math/big"
	"net/url"

	"github.com/google/uuid"
	"github.com/sqids/sqids-go"

	"github.com/patrick-devel/shorturl/internal/models"
)

const minLength = 6

type ShortLinkService struct {
	baseURL *url.URL
	storage storage
}

func New(baseURL *url.URL, storage storage) *ShortLinkService {
	return &ShortLinkService{baseURL: baseURL, storage: storage}
}

type storage interface {
	ReadEvent(ctx context.Context, hash string) (string, error)
	WriteEvent(ctx context.Context, event models.Event) error
	WriteEvents(ctx context.Context, event []models.Event) error
}

func (sh *ShortLinkService) MakeShortURL(ctx context.Context, originalURL, uid string) (string, error) {
	hash, err := sh.generateHash(originalURL)
	if err != nil {
		return "", fmt.Errorf("genarate hash failed: %w", err)
	}

	if uid == "" {
		uid = uuid.NewString()
	}

	event := models.Event{
		UUID:        uid,
		ShortURL:    sh.baseURL.String() + "/" + hash,
		OriginalURL: originalURL,
		Hash:        hash,
	}

	err = sh.storage.WriteEvent(ctx, event)
	if err != nil {
		return "", fmt.Errorf("save event failed: %w", err)
	}
	return event.ShortURL, nil
}

func (sh *ShortLinkService) GetOriginalURL(ctx context.Context, hash string) (string, error) {
	originalURL, err := sh.storage.ReadEvent(ctx, hash)
	if err != nil {
		return "", fmt.Errorf("fetch url failed or not found: %w", err)
	}

	return originalURL, nil
}

func (sh *ShortLinkService) generateHash(url string) (string, error) {
	generatedNumber := new(big.Int).SetBytes([]byte(url)).Uint64()
	s, err := sqids.New(sqids.Options{MinLength: minLength})
	if err != nil {
		return "", err
	}

	id, err := s.Encode([]uint64{generatedNumber})
	if err != nil {
		return "", err
	}

	return id, nil
}

func (sh *ShortLinkService) MakeShortURLs(ctx context.Context, bulk models.ListRequestBulk) ([]models.Event, error) {
	events := make([]models.Event, 0, len(bulk))
	for _, r := range bulk {
		hash, err := sh.generateHash(r.OriginalURL.String())
		if err != nil {
			return events, fmt.Errorf("genarate hash failed: %w", err)
		}

		event := models.Event{
			UUID:        r.CorrelationID,
			ShortURL:    sh.baseURL.String() + "/" + hash,
			OriginalURL: r.OriginalURL.String(),
			Hash:        hash,
		}
		events = append(events, event)
	}

	err := sh.storage.WriteEvents(ctx, events)
	if err != nil {
		return events, fmt.Errorf("save events failed: %w", err)
	}

	return events, nil
}
