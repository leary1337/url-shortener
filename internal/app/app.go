package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/leary1337/url-shortener/internal/app/config"
	"github.com/leary1337/url-shortener/internal/app/handler"
	"github.com/leary1337/url-shortener/internal/app/middleware"
	"github.com/leary1337/url-shortener/internal/app/service"
	"github.com/leary1337/url-shortener/internal/app/service/repo"
	"github.com/leary1337/url-shortener/pkg/logger"
)

func RunServer(cfg *config.Config) error {
	l := logger.New(cfg.Log.Level)

	shortenerRepo := repo.NewShortenerMemory()
	shortenerSrv := service.NewShortenerService(shortenerRepo)
	shortenerHandler := handler.NewShortenerHandler(l, shortenerSrv, cfg.RedirectAddr)

	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware(l))
	r.Route("/", func(r chi.Router) {
		shortenerHandler.RegisterRoutes(r)
	})
	return http.ListenAndServe(cfg.Addr, r)
}
