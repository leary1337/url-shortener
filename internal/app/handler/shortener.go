package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/leary1337/url-shortener/internal/app/entity"
	"github.com/leary1337/url-shortener/internal/app/service"
	"github.com/leary1337/url-shortener/pkg/logger"
)

type ShortenerHandler struct {
	l            logger.Interface
	service      service.Shortener
	redirectAddr string
}

func NewShortenerHandler(l logger.Interface, service service.Shortener, redirectAddr string) *ShortenerHandler {
	return &ShortenerHandler{
		l:            l,
		service:      service,
		redirectAddr: redirectAddr,
	}
}
func (s *ShortenerHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", s.ShortenURL)
	r.Get("/{shortURL}", s.ResolveURL)
	r.Post("/api/shorten", s.ShortenURLJSON)
}

func (s *ShortenerHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	originalURL := strings.TrimSpace(string(body))
	shortURL, err := s.service.ShortenURL(originalURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL.GetFullShortURL(s.redirectAddr)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *ShortenerHandler) ShortenURLJSON(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var srb entity.ShortenRequestBody
	err = json.Unmarshal(body, &srb)
	if err != nil || srb.Url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := s.service.ShortenURL(srb.Url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseBody, err := json.Marshal(&entity.ShortenResponseBody{Result: shortURL.GetFullShortURL(s.redirectAddr)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *ShortenerHandler) ResolveURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	url, err := s.service.ResolveURL(shortURL)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Location", url.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
