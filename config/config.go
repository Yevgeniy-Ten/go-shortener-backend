package config

import (
	"flag"
	handlers "shorter/internal/handlers"
)

type Config struct {
	Address *string
	Config  handlers.Config
}

func NewConfig() *Config {
	address := flag.String("a", ":8080", "address for server")
	serveAddr := flag.String("b", "http://localhost:8080", "address for link")
	flag.Parse()

	config := Config{
		Address: address,
		Config:  handlers.Config{ServerAddr: *serveAddr},
	}

	return &config
}
