package service

import (
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

func (s *ShortenerService) ShortenURL(originalURL string) string {
	shortURL := util.GenerateShortURL(ShortURLLength)
	s.repo.Save(shortURL, originalURL)
	return shortURL
}

func (s *ShortenerService) ResolveURL(shortURL string) (string, error) {
	originalURL, exists := s.repo.Get(shortURL)
	if !exists {
		return "", ErrURLNotFound
	}
	return originalURL, nil
}
