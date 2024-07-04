package config

import (
	"flag"
	"os"
)

type (
	Config struct {
		Log
		Addr         string `env:"SERVER_ADDRESS"`
		RedirectAddr string `env:"BASE_URL"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL"`
	}
)

const (
	DefaultServerAddr   = `localhost:8080`
	DefaultRedirectAddr = `http://localhost:8080`
)

func NewConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Addr, "a", DefaultServerAddr, "server address")
	flag.StringVar(&cfg.RedirectAddr, "b", DefaultRedirectAddr, "redirect server address")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		cfg.Addr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		cfg.RedirectAddr = envBaseAddr
	}
	return &cfg, nil
}
