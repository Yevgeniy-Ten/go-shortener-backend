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
	port := flag.String("a", "8080", "port for server")
	serveAddr := flag.String("b", "", "address for link")
	flag.Parse()

	if *serveAddr == "" {
		*serveAddr = "http://localhost:" + *port + "/"
	}

	config := Config{
		Port:   port,
		Config: handlers.Config{ServerAddr: *serveAddr},
	}

	return &config
}
