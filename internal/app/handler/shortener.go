package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/leary1337/url-shortener/internal/app/service"
)

type ShortenerHandler struct {
	service      service.Shortener
	redirectAddr string
}

func NewShortenerHandler(service service.Shortener, redirectAddr string) *ShortenerHandler {
	return &ShortenerHandler{
		service:      service,
		redirectAddr: redirectAddr,
	}
}
func (h *ShortenerHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.ShortenURL)
	r.Get("/{shortURL}", h.ResolveURL)
}

func (h *ShortenerHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	originalURL := strings.TrimSpace(string(body))
	shortURL := h.service.ShortenURL(originalURL)

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("%s/%s", h.redirectAddr, shortURL)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *ShortenerHandler) ResolveURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	originalURL, err := h.service.ResolveURL(shortURL)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
