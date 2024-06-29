package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ShortenerRouter(sh *ServerHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get(`/{shortURL}`, sh.GetOriginalURL)
	r.Post(`/`, sh.GenerateShortURL)
	return r
}
