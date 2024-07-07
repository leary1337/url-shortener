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
	l       logger.Interface
	service service.Shortener
}

func NewShortenerHandler(l logger.Interface, service service.Shortener) *ShortenerHandler {
	return &ShortenerHandler{
		l:       l,
		service: service,
	}
}
func (s *ShortenerHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", s.ShortenURL)
	r.Get("/{shortURL}", s.ResolveURL)
	r.Post("/api/shorten", s.ShortenURLJSON)
	r.Post("/api/shorten/batch", s.ShortenURLBatch)
}

func (s *ShortenerHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	originalURL := strings.TrimSpace(string(body))
	shortURL, err := s.service.ShortenURL(r.Context(), originalURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
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

	shortURL, err := s.service.ShortenURL(r.Context(), srb.Url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseBody, err := json.Marshal(&entity.ShortenResponseBody{Result: shortURL})
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

func (s *ShortenerHandler) ShortenURLBatch(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var sb []entity.ShortenBatchRequestBody
	err = json.Unmarshal(body, &sb)
	if err != nil || len(sb) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := s.service.ShortenBatch(r.Context(), sb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(resp)
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
	originalURL, err := s.service.ResolveURL(r.Context(), shortURL)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
