package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/leary1337/url-shortener/internal/app"
	"github.com/leary1337/url-shortener/internal/app/config"
)

func main() {
	c := config.NewConfig()
	fmt.Printf("Config: %+v\n", c)
	serverHandler := app.NewServerHandler(c)
	err := http.ListenAndServe(c.Addr, app.ShortenerRouter(serverHandler))
	if err != nil {
		log.Fatal(err)
	}
}
