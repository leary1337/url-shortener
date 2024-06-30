package main

import (
	"log"
	"net/http"

	"github.com/leary1337/url-shortener/internal/app"
	"github.com/leary1337/url-shortener/internal/app/config"
)

func main() {
	c := config.NewConfig()
	serverHandler := app.NewServerHandler(c)
	err := http.ListenAndServe(c.Addr, app.ShortenerRouter(serverHandler))
	if err != nil {
		log.Fatal(err)
	}
}
