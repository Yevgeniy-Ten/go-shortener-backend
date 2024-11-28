package config

import (
	"flag"
	handlers "shorter/internal/handlers"
)

type Config struct {
	Port   *string
	Config handlers.Config
}

func NewConfig() *Config {
	port := flag.String("a", "localhost:8080", "port for server")
	serveAddr := flag.String("b", "http://localhost:8080", "address for link")
	flag.Parse()

	config := Config{
		Port:   port,
		Config: handlers.Config{ServerAddr: *serveAddr},
	}

	return &config
}
