package repo

import (
	"encoding/json"
	"os"

	"github.com/leary1337/url-shortener/internal/app/entity"
	"github.com/leary1337/url-shortener/internal/app/service"
)

var _ service.ShortenerRepo = (*ShortenerFileStorage)(nil)

type ShortenerFileStorage struct {
	filePath string
	m        *ShortenerMemory
}

func NewShortenerFileStorage(filePath string) *ShortenerFileStorage {
	s := &ShortenerFileStorage{
		filePath: filePath,
		m:        NewShortenerMemory(),
	}
	s.loadToMemory()
	return s
}

func (s *ShortenerFileStorage) Save(shortURL *entity.ShortURL) error {
	// Save to memory
	err := s.m.Save(shortURL)
	if err != nil {
		return err
	}

	data, err := json.Marshal(s.m.GetAll())
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0666)
}

func (s *ShortenerFileStorage) GetByShortURL(shortURL string) (*entity.ShortURL, error) {
	return s.m.GetByShortURL(shortURL)
}

func (s *ShortenerFileStorage) loadToMemory() {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return
	}

	var urls []entity.ShortURL
	if err = json.Unmarshal(data, &urls); err != nil {
		return
	}

	for _, url := range urls {
		_ = s.m.Save(&url)
	}
}
