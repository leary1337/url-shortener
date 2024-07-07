package repo

import (
	"context"
	"sync"

	"github.com/leary1337/url-shortener/internal/app/entity"
	"github.com/leary1337/url-shortener/internal/app/service"
)

var _ service.ShortenerRepo = (*ShortenerMemory)(nil)

type ShortenerMemory struct {
	mu    sync.RWMutex
	store map[string]entity.ShortURL
}

func NewShortenerMemory() *ShortenerMemory {
	return &ShortenerMemory{
		store: make(map[string]entity.ShortURL),
	}
}

func (m *ShortenerMemory) Save(ctx context.Context, shortURL *entity.ShortURL) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[shortURL.ShortURI] = *shortURL
	return nil
}

func (m *ShortenerMemory) SaveBatch(ctx context.Context, shortURLs []entity.ShortURL) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, url := range shortURLs {
		m.store[url.ShortURI] = url
	}
	return nil
}

func (m *ShortenerMemory) GetByShortURI(ctx context.Context, shortURL string) (*entity.ShortURL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	url, ok := m.store[shortURL]
	if !ok {
		return nil, service.ErrURLNotFound
	}
	return &url, nil
}

func (m *ShortenerMemory) GetAll() []entity.ShortURL {
	urls := make([]entity.ShortURL, 0, len(m.store))
	for _, url := range m.store {
		urls = append(urls, url)
	}
	return urls
}
