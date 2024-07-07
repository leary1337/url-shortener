package repo

import (
	"context"
	"encoding/json"
	"os"

	"github.com/leary1337/url-shortener/internal/app/entity"
	"github.com/leary1337/url-shortener/internal/app/service"
)

var _ service.ShortenerRepo = (*ShortenerFileMemory)(nil)

type ShortenerFileMemory struct {
	filePath string
	m        *ShortenerMemory
}

func NewShortenerFileMemory(filePath string) *ShortenerFileMemory {
	s := &ShortenerFileMemory{
		filePath: filePath,
		m:        NewShortenerMemory(),
	}
	s.loadToMemory()
	return s
}

func (s *ShortenerFileMemory) Save(ctx context.Context, shortURL *entity.ShortURL) error {
	// Save to memory
	err := s.m.Save(ctx, shortURL)
	if err != nil {
		return err
	}

	data, err := json.Marshal(s.m.GetAll())
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0666)
}

func (s *ShortenerFileMemory) GetByShortURL(ctx context.Context, shortURL string) (*entity.ShortURL, error) {
	return s.m.GetByShortURL(ctx, shortURL)
}

func (s *ShortenerFileMemory) loadToMemory() {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return
	}

	var urls []entity.ShortURL
	if err = json.Unmarshal(data, &urls); err != nil {
		return
	}

	for _, url := range urls {
		_ = s.m.Save(context.Background(), &url)
	}
}
