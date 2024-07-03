package main

import (
	"log"
	"net/http"

	"github.com/leary1337/url-shortener/internal/app"
	"github.com/leary1337/url-shortener/internal/app/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	serverHandler := app.NewServerHandler(cfg)

	err = http.ListenAndServe(cfg.Addr, app.ShortenerRouter(serverHandler))
	if err != nil {
		log.Fatal(err)
	}
}
