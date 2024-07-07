package service

import (
	"context"

	"github.com/leary1337/url-shortener/internal/app/entity"
)

type (
	Shortener interface {
		ShortenURL(originalURL string) (*entity.ShortURL, error)
		ResolveURL(shortURL string) (*entity.ShortURL, error)
	}

	ShortenerRepo interface {
		Save(shortURL *entity.ShortURL) error
		GetByShortURL(shortURL string) (*entity.ShortURL, error)
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
