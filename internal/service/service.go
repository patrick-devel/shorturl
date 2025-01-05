package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/sqids/sqids-go"

	"github.com/patrick-devel/shorturl/internal/ctxaux"
	"github.com/patrick-devel/shorturl/internal/models"
	"github.com/patrick-devel/shorturl/internal/storage"
)

const minLength = 6
const batchDelete = 10

type ShortLinkService struct {
	baseURL *url.URL
	storage store

	urlsCh chan string
	ctx    context.Context
}

func New(baseURL *url.URL, storage store, ctx context.Context) *ShortLinkService {
	urlsCh := make(chan string)
	sh := &ShortLinkService{baseURL: baseURL, storage: storage, urlsCh: urlsCh, ctx: ctx}
	go sh.runDelete()
	return sh
}

type store interface {
	ReadEvent(ctx context.Context, hash string) (string, error)
	WriteEvent(ctx context.Context, event models.Event) error
	WriteEvents(_ context.Context, events []models.Event) error
	ReadEventsByCreatorID(ctx context.Context, userID string) ([]models.Event, error)
	SetDeleteByShortURL(shorts []string) error
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

func (sh *ShortLinkService) DeleteShortURL(ctx context.Context, shortUrls []string) error {
	if ctxaux.GetUserIDFromContext(ctx) == "" {
		return errors.New("no user is currently logged in")
	}

	chUrls := sh.urlDeleteGenerator(ctx, shortUrls)
	linksByUser, err := sh.LinksByCreatorID(ctx)
	if err != nil {
		return fmt.Errorf("get links by creator id failed: %w", err)
	}

	URLForDeleteCh := sh.sendDeleteByUser(ctx, chUrls, linksByUser)
	for v := range URLForDeleteCh {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sh.urlsCh <- v:
		}
	}

	return nil
}

func (sh *ShortLinkService) urlDeleteGenerator(ctx context.Context, shortUrls []string) chan string {
	checkCh := make(chan string)
	go func() {
		defer close(checkCh)

		for _, u := range shortUrls {
			select {
			case checkCh <- sh.baseURL.String() + "/" + u:
			case <-ctx.Done():
				return
			}
		}
	}()

	return checkCh
}

func (sh *ShortLinkService) sendDeleteByUser(ctx context.Context, urls chan string, urlsByUser []models.Event) chan string {
	resURL := make(chan string)
	go func() {
		defer close(resURL)

		for {
			select {
			case data := <-urls:
				for _, u := range urlsByUser {
					if u.ShortURL == data {
						resURL <- data
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return resURL
}

func (sh *ShortLinkService) runDelete() {
	urlsBatch := make([]string, 0, batchDelete*2)
	tiker := time.NewTicker(1 * time.Second)

	funcDelete := func() {
		if len(urlsBatch) != 0 {
			err := sh.storage.SetDeleteByShortURL(urlsBatch)
			if err != nil {
				logrus.Errorf("set delete batch failed: %v", err)
			} else {
				urlsBatch = []string{}
			}
		}
	}

	for {
		select {
		case data := <-sh.urlsCh:
			urlsBatch = append(urlsBatch, data)
			if len(urlsBatch) >= batchDelete {
				funcDelete()
			}
		case <-tiker.C:
			funcDelete()
		case <-sh.ctx.Done():
			funcDelete()
			return
		}
	}
}
