package entity

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type ShortenRequestBody struct {
	Url string `json:"url"`
}

type ShortenResponseBody struct {
	Result string `json:"result"`
}

type ShortURL struct {
	UUID        uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

func (s *ShortURL) GetFullShortURL(addr string) string {
	return fmt.Sprintf("%s/%s", addr, s.ShortURL)
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
