package entity

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/leary1337/url-shortener/internal/app/util"
)

type (
	ShortenRequestBody struct {
		Url string `json:"url"`
	}
	ShortenResponseBody struct {
		Result string `json:"result"`
	}
)

type (
	ShortenBatchRequestBody struct {
		Id          string `json:"correlation_id"`
		OriginalURL string `json:"original_url"`
	}
	ShortenBatchResponseBody struct {
		Id       string `json:"correlation_id"`
		ShortURL string `json:"short_url"`
	}
)

type ShortURL struct {
	UUID        uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

const ShortURLLength = 8

func NewShortURL(originalURL, addr string) *ShortURL {
	return &ShortURL{
		UUID:        uuid.New(),
		ShortURL:    util.GenerateShortURL(addr, ShortURLLength),
		OriginalURL: originalURL,
	}
}

func (s *ShortURL) MarshalJSON() ([]byte, error) {
	type Alias ShortURL // Создаем алиас для структуры ShortURL
	return json.Marshal(&struct {
		UUID string `json:"uuid"`
		*Alias
	}{
		UUID:  s.UUID.String(),
		Alias: (*Alias)(s), // Используем алиас для встраивания всех полей ShortURL
	})
}

func (s *ShortURL) UnmarshalJSON(data []byte) error {
	type Alias ShortURL // Создаем алиас для структуры ShortURL
	aux := &struct {
		UUID string `json:"uuid"`
		*Alias
	}{
		Alias: (*Alias)(s), // Используем алиас для встраивания всех полей ShortURL
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	var err error
	s.UUID, err = uuid.Parse(aux.UUID) // Преобразование строки обратно в UUID
	return err
}
