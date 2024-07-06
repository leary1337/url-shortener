package service

import "github.com/leary1337/url-shortener/internal/app/entity"

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
