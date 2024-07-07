package entity

import (
	"encoding/json"
	"fmt"

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
	ShortURI    string    `json:"short_uri"`
	OriginalURL string    `json:"original_url"`
}

const ShortURLLength = 8

func NewShortURL(originalURL string) *ShortURL {
	return &ShortURL{
		UUID:        uuid.New(),
		ShortURI:    util.RandomString(ShortURLLength),
		OriginalURL: originalURL,
	}
}

func (s *ShortURL) MarshalJSON() ([]byte, error) {
	type Alias ShortURL // Создаем алиас для структуры ShortURI
	return json.Marshal(&struct {
		UUID string `json:"uuid"`
		*Alias
	}{
		UUID:  s.UUID.String(),
		Alias: (*Alias)(s), // Используем алиас для встраивания всех полей ShortURI
	})
}

func (s *ShortURL) UnmarshalJSON(data []byte) error {
	type Alias ShortURL // Создаем алиас для структуры ShortURI
	aux := &struct {
		UUID string `json:"uuid"`
		*Alias
	}{
		Alias: (*Alias)(s), // Используем алиас для встраивания всех полей ShortURI
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	var err error
	s.UUID, err = uuid.Parse(aux.UUID) // Преобразование строки обратно в UUID
	return err
}

func (s *ShortURL) GetShortURL(addr string) string {
	return fmt.Sprintf("%s/%s", addr, s.ShortURI)
}
