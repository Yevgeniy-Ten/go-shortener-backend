package config

import (
	"flag"
	"fmt"
	handlers "shorter/internal/handlers"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address string
	Config  *handlers.Config
}

func NewConfig() (*Config, error) {
	var config = Config{
		Config: &handlers.Config{
			ServerAddr: "",
		},
		Address: "",
	}
	parseFlags(&config)
	err := parseEnv(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func parseEnv(config *Config) error {
	type EnvConfig struct {
		Address    string `env:"SERVER_ADDRESS"`
		ServerAddr string `env:"SERVER_URL"`
	}

	var cfg EnvConfig
	if err := env.Parse(&cfg); err != nil {
		return err
	}
	if cfg.Address != "" {
		config.Address = cfg.Address
		fmt.Println("Address from env", config.Address)
	}
	if cfg.ServerAddr != "" {
		config.Config.ServerAddr = cfg.ServerAddr
	}
	return nil
}
func parseFlags(config *Config) {
	address := flag.String("a", ":8080", "address for server")
	serverAddr := flag.String("b", "http://localhost:8080", "address for link")
	flag.Parse()
	config.Address = *address
	config.Config.ServerAddr = *serverAddr
}
