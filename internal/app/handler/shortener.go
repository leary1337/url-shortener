package handler

import (
	"encoding/json"
	"fmt"
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
func (h *ShortenerHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.ShortenURL)
	r.Get("/{shortURL}", h.ResolveURL)
	r.Post("/api/shorten", h.ShortenURLJSON)
}

func (h *ShortenerHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	originalURL := strings.TrimSpace(string(body))
	shortURL := h.service.ShortenURL(originalURL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("%s/%s", h.redirectAddr, shortURL)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *ShortenerHandler) ShortenURLJSON(w http.ResponseWriter, r *http.Request) {
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

	shortURL := fmt.Sprintf("%s/%s", h.redirectAddr, h.service.ShortenURL(srb.Url))

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
