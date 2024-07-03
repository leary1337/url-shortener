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

	// Сначала определите все флаги.
	flag.StringVar(&cfg.Addr, "a", "", "server address")
	flag.StringVar(&cfg.RedirectAddr, "b", "", "redirect server address")

	// Парсинг флагов перед парсингом переменных окружения, чтобы установить пустые значения.
	flag.Parse()

	// Парсинг переменных окружения.
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	// Применение значений по умолчанию, если не заданы ни флаги, ни переменные окружения.
	if cfg.Addr == "" {
		cfg.Addr = DefaultServerAddr
	}
	if cfg.RedirectAddr == "" {
		cfg.RedirectAddr = DefaultRedirectAddr
	}

	return &cfg, nil
}
