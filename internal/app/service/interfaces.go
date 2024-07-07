package service

import (
	"context"

	"github.com/leary1337/url-shortener/internal/app/entity"
)

type (
	Shortener interface {
		ShortenURL(ctx context.Context, originalURL string) (string, error)
		ShortenBatch(ctx context.Context, sb []entity.ShortenBatchRequestBody) ([]entity.ShortenBatchResponseBody, error)
		ResolveURL(ctx context.Context, shortURL string) (string, error)
	}

	ShortenerRepo interface {
		Save(ctx context.Context, shortURL *entity.ShortURL) error
		SaveBatch(ctx context.Context, shortURLs []entity.ShortURL) error
		GetByShortURI(ctx context.Context, shortURL string) (*entity.ShortURL, error)
	}
)

type (
	Ping interface {
		PingDB(ctx context.Context) error
	}

	PingRepo interface {
		Ping(ctx context.Context) error
	}
)
