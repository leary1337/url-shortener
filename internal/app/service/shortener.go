package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/leary1337/url-shortener/internal/app/entity"
	"github.com/leary1337/url-shortener/internal/app/util"
)

var _ Shortener = (*ShortenerService)(nil)

const ShortURLLength = 8

type ShortenerService struct {
	repo ShortenerRepo
}

func NewShortenerService(repo ShortenerRepo) *ShortenerService {
	return &ShortenerService{repo: repo}
}

func (s *ShortenerService) ShortenURL(ctx context.Context, originalURL string) (*entity.ShortURL, error) {
	shortURL := util.GenerateShortURL(ShortURLLength)
	url := &entity.ShortURL{
		UUID:        uuid.New(),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
	err := s.repo.Save(ctx, url)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (s *ShortenerService) ResolveURL(ctx context.Context, shortURL string) (*entity.ShortURL, error) {
	url, err := s.repo.GetByShortURL(ctx, shortURL)
	if err != nil {
		return nil, err
	}
	return url, nil
}
