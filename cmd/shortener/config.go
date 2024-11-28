package main

import (
	"flag"
	handlers "shorter/internal/handlers"
)

type Config struct {
	port   *string
	Config handlers.Config
}

func NewConfig() *Config {
	serveAddr := flag.String("b", "http://localhost:8080/", "port for server")
	config := Config{
		port:   flag.String("a", "8080", "port for server"),
		Config: handlers.Config{ServerAddr: *serveAddr},
	}

	return &config
}
