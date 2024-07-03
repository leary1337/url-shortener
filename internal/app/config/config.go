package config

import (
	"flag"
	"github.com/caarlos0/env"
)

type Config struct {
	Addr         string `env:"SERVER_ADDRESS"`
	RedirectAddr string `env:"BASE_URL"`
}

const (
	DefaultServerAddr   = `localhost:8080`
	DefaultRedirectAddr = `http://localhost:8080`
)

func NewConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Addr == "" {
		flag.StringVar(&cfg.Addr, "a", DefaultServerAddr, "server address")
	}
	if cfg.RedirectAddr == "" {
		flag.StringVar(&cfg.RedirectAddr, "b", DefaultRedirectAddr, "redirect server address")
	}
	flag.Parse()
	return &cfg, nil
}
