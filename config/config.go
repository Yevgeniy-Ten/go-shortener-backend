package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	handlers "shorter/internal/handlers"
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
	}
	if cfg.ServerAddr != "" {
		config.Config.ServerAddr = cfg.ServerAddr
	}
	return nil
}
func parseFlags(config *Config) {
	config.Address = *flag.String("a", ":8080", "address for server")
	config.Config.ServerAddr = *flag.String("b", "http://localhost:8080", "address for link")
	flag.Parse()
}
