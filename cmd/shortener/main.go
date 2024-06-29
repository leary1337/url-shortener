package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/leary1337/url-shortener/internal/app"
)

const ServerAddress = `localhost:8080`

func main() {
	serverHandler := app.NewServerHandler(ServerAddress)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get(`/{shortURL}`, serverHandler.GetOriginalURL)
	r.Post(`/`, serverHandler.GenerateShortURL)

	err := http.ListenAndServe(ServerAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}
