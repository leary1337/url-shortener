package app

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

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

	// Init Postgres db pool
	dbPool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		return err
	}
	defer dbPool.Close()

	pingRepo := repo.NewPingPostgres(dbPool)
	pingSrv := service.NewPingService(pingRepo)
	pingHandler := handler.NewPingHandler(l, pingSrv)

	//shortenerRepo := repo.NewShortenerMemory()
	shortenerRepo := repo.NewShortenerFileMemory(cfg.FileStoragePath)
	shortenerSrv := service.NewShortenerService(shortenerRepo)
	shortenerHandler := handler.NewShortenerHandler(l, shortenerSrv, cfg.RedirectAddr)

	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware(l))
	r.Use(middleware.CompressMiddleware)
	r.Route("/", func(r chi.Router) {
		pingHandler.RegisterRoutes(r)
		shortenerHandler.RegisterRoutes(r)
	})
	return http.ListenAndServe(cfg.Addr, r)
}
