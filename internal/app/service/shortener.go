package service

import (
	"context"

	"github.com/leary1337/url-shortener/internal/app/entity"
)

var _ Shortener = (*ShortenerService)(nil)

type ShortenerService struct {
	repo         ShortenerRepo
	redirectAddr string
}

func NewShortenerService(repo ShortenerRepo, addr string) *ShortenerService {
	return &ShortenerService{repo: repo, redirectAddr: addr}
}

func (s *ShortenerService) ShortenURL(ctx context.Context, originalURL string) (*entity.ShortURL, error) {
	url := entity.NewShortURL(originalURL, s.redirectAddr)
	err := s.repo.Save(ctx, url)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (s *ShortenerService) ShortenBatch(ctx context.Context, sb []entity.ShortenBatchRequestBody) ([]entity.ShortenBatchResponseBody, error) {
	urls := make([]entity.ShortURL, 0, len(sb))
	rb := make([]entity.ShortenBatchResponseBody, 0, len(sb))
	for _, r := range sb {
		url := entity.NewShortURL(r.OriginalURL, s.redirectAddr)
		urls = append(urls, *url)
		rb = append(rb, entity.ShortenBatchResponseBody{
			Id:       r.Id,
			ShortURL: url.ShortURL,
		})
	}
	err := s.repo.SaveBatch(ctx, urls)
	if err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *ShortenerService) ResolveURL(ctx context.Context, shortURL string) (*entity.ShortURL, error) {
	url, err := s.repo.GetByShortURL(ctx, shortURL)
	if err != nil {
		return nil, err
	}
	return url, nil
}
