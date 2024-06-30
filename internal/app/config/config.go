package config

import "flag"

type Config struct {
	Addr         string
	RedirectAddr string
}

const DefaultServerAddr = `localhost:8080`

func NewConfig() *Config {
	config := &Config{}
	flag.StringVar(&config.Addr, "a", DefaultServerAddr, "server address")
	flag.StringVar(&config.RedirectAddr, "b", DefaultServerAddr, "redirect server address")
	flag.Parse()
	return config
}
