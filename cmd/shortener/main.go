package main

import (
	"log"
	"net/http"

	"github.com/leary1337/url-shortener/internal/app"
)

const ServerAddress = `localhost:8080`

func main() {
	serverHandler := app.NewServerHandler(ServerAddress)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, serverHandler.MainPage)

	err := http.ListenAndServe(ServerAddress, mux)
	if err != nil {
		log.Fatal(err)
	}
}
