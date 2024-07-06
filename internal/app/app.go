package app

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/leary1337/url-shortener/internal/app/config"
	"github.com/leary1337/url-shortener/internal/app/handler"
	"github.com/leary1337/url-shortener/internal/app/middleware"
	"github.com/leary1337/url-shortener/internal/app/service"
	"github.com/leary1337/url-shortener/internal/app/service/repo"
	"github.com/leary1337/url-shortener/pkg/logger"
)

func RunServer(cfg *config.Config) error {
	//logFile, err := os.OpenFile("shortener.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	//if err != nil {
	//	return err
	//}
	//defer func() {
	//	_ = logFile.Close()
	//}()
	l := logger.New(cfg.Log.Level, os.Stdout)

	shortenerRepo := repo.NewShortenerMemory()
	shortenerSrv := service.NewShortenerService(shortenerRepo)
	shortenerHandler := handler.NewShortenerHandler(l, shortenerSrv, cfg.RedirectAddr)

	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware(l))
	r.Use(middleware.CompressMiddleware)
	r.Route("/", func(r chi.Router) {
		shortenerHandler.RegisterRoutes(r)
	})
	return http.ListenAndServe(cfg.Addr, r)
}
