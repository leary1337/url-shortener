package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/leary1337/url-shortener/internal/app/config"
	"github.com/leary1337/url-shortener/internal/app/handler"
	"github.com/leary1337/url-shortener/internal/app/service"
	"github.com/leary1337/url-shortener/internal/app/service/repo"
)

func RunServer(cfg *config.Config) error {
	shortenerRepo := repo.NewShortenerMemory()
	shortenerSrv := service.NewShortenerService(shortenerRepo)
	shortenerHandler := handler.NewShortenerHandler(shortenerSrv, cfg.RedirectAddr)

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		shortenerHandler.RegisterRoutes(r)
	})
	return http.ListenAndServe(cfg.Addr, r)
}
