package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"

	"github.com/google/uuid"
	"github.com/sqids/sqids-go"

	"github.com/patrick-devel/shorturl/internal/ctxaux"
	"github.com/patrick-devel/shorturl/internal/models"
	"github.com/patrick-devel/shorturl/internal/storage"
)

const minLength = 6

type ShortLinkService struct {
	baseURL *url.URL
	storage istorage
}

func New(baseURL *url.URL, storage istorage) *ShortLinkService {
	return &ShortLinkService{baseURL: baseURL, storage: storage}
}

type istorage interface {
	ReadEvent(ctx context.Context, hash string) (string, error)
	WriteEvent(ctx context.Context, event models.Event) error
	WriteEvents(ctx context.Context, event []models.Event) error
	ReadEventsByCreatorID(ctx context.Context, userID string) ([]models.Event, error)
}

func (sh *ShortLinkService) MakeShortURL(ctx context.Context, originalURL, uid string) (string, error) {
	hash, err := sh.generateHash(originalURL)
	if err != nil {
		return "", fmt.Errorf("generate hash failed: %w", err)
	}

	if uid == "" {
		uid = uuid.NewString()
	}

	event := models.Event{
		UUID:        uid,
		CreatorID:   ctxaux.GetUserIDFromContext(ctx),
		ShortURL:    sh.baseURL.String() + "/" + hash,
		OriginalURL: originalURL,
	}

	err = sh.storage.WriteEvent(ctx, event)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateURL) {
			return event.ShortURL, err
		}

		return "", fmt.Errorf("save event failed: %w", err)
	}
	return event.ShortURL, nil
}

func (sh *ShortLinkService) GetOriginalURL(ctx context.Context, hash string) (string, error) {
	short := sh.baseURL.String() + hash
	originalURL, err := sh.storage.ReadEvent(ctx, short)
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
			CreatorID:   ctxaux.GetUserIDFromContext(ctx),
			ShortURL:    sh.baseURL.String() + "/" + hash,
			OriginalURL: r.OriginalURL.String(),
		}
		events = append(events, event)
	}

	err := sh.storage.WriteEvents(ctx, events)
	if err != nil {
		return events, fmt.Errorf("save events failed: %w", err)
	}

	return events, nil
}

func (sh *ShortLinkService) LinksByCreatorID(ctx context.Context) ([]models.Event, error) {
	var events []models.Event

	events, err := sh.storage.ReadEventsByCreatorID(ctx, ctxaux.GetUserIDFromContext(ctx))
	if err != nil {
		return events, fmt.Errorf("failed get links for current user")
	}

	return events, nil
}
