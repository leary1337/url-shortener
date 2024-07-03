package main

import (
	"log"
	"net/http"

	"github.com/leary1337/url-shortener/internal/app"
	"github.com/leary1337/url-shortener/internal/app/config"
)

func main() {
	c, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	serverHandler := app.NewServerHandler(c)

	err = http.ListenAndServe(c.Addr, app.ShortenerRouter(serverHandler))
	if err != nil {
		log.Fatal(err)
	}
}
