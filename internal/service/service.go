package service

import (
	"context"
	"fmt"
	"math/big"
	"net/url"

	"github.com/sqids/sqids-go"
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
	WriteEvent(ctx context.Context, hash, originalURL string) error
}

func (sh *ShortLinkService) MakeShortURL(ctx context.Context, originalURL string) (string, error) {
	hash, err := sh.generateHash(originalURL)
	if err != nil {
		return "", fmt.Errorf("genarate hash failed: %w", err)
	}

	err = sh.storage.WriteEvent(ctx, hash, originalURL)
	if err != nil {
		return "", fmt.Errorf("save event failed: %w", err)
	}
	return sh.baseURL.String() + "/" + hash, nil
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
