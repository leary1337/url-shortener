package main

import (
	"log"
	"net/http"

	"github.com/leary1337/url-shortener/internal/app"
)

const ServerAddress = `localhost:8080`

func main() {
	serverHandler := app.NewServerHandler(ServerAddress)
	err := http.ListenAndServe(ServerAddress, app.ShortenerRouter(serverHandler))
	if err != nil {
		log.Fatal(err)
	}
}
