package main

import (
	"log"

	"github.com/leary1337/url-shortener/internal/app"
	"github.com/leary1337/url-shortener/internal/app/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	if err = app.RunServer(cfg); err != nil {
		log.Fatalln(err)
	}
}
