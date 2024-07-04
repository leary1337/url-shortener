package repo

import (
	"sync"

	"github.com/leary1337/url-shortener/internal/app/service"
)

var _ service.ShortenerRepo = (*ShortenerMemory)(nil)

type ShortenerMemory struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewShortenerMemory() *ShortenerMemory {
	return &ShortenerMemory{
		store: make(map[string]string),
	}
}

func (m *ShortenerMemory) Save(shortURL, originalURL string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[shortURL] = originalURL
}

func (m *ShortenerMemory) Get(shortURL string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	originalURL, exists := m.store[shortURL]
	return originalURL, exists
}
