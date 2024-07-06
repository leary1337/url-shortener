package config

import (
	"flag"
	"os"
)

type (
	Config struct {
		Log
		DB
		Addr         string `env:"SERVER_ADDRESS"`
		RedirectAddr string `env:"BASE_URL"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL"`
	}

	DB struct {
		FileStoragePath string `env:"FILE_STORAGE_PATH"`
	}
)

const (
	DefaultServerAddr        = `localhost:8080`
	DefaultRedirectAddr      = `http://localhost:8080`
	DefaultDBFileStoragePath = `/tmp/short-url-db.json`
)

func NewConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Addr, "a", DefaultServerAddr, "server address")
	flag.StringVar(&cfg.RedirectAddr, "b", DefaultRedirectAddr, "redirect server address")
	flag.StringVar(&cfg.FileStoragePath, "f", DefaultDBFileStoragePath, "file storage path")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		cfg.Addr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		cfg.RedirectAddr = envBaseAddr
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}
	return &cfg, nil
}
