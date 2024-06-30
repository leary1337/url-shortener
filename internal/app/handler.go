package app

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/leary1337/url-shortener/internal/app/config"
)

const ShortURLLength = 8

type ServerHandler struct {
	cfg    *config.Config
	urlMap map[string]string
}

func NewServerHandler(cfg *config.Config) *ServerHandler {
	return &ServerHandler{
		cfg:    cfg,
		urlMap: map[string]string{},
	}
}

func (a *ServerHandler) GenerateShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := GenerateRandomString(ShortURLLength)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	a.urlMap[shortURL] = string(body)

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("%s/%s", a.cfg.RedirectAddr, shortURL)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *ServerHandler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	originalURL, ok := a.urlMap[shortURL]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
