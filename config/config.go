package config

import (
	"flag"
	"shorter/internal/handlers"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address  string
	FilePath string
	Config   *handlers.Config
}

func NewConfig() (*Config, error) {
	config := &Config{
		Address:  ":8080",
		FilePath: "storage",
		Config: &handlers.Config{
			ServerAddr: "http://localhost:8080",
		},
	}

	parseFlags(config)
	if err := parseEnv(config); err != nil {
		return nil, err
	}

	return config, nil
}

func parseEnv(config *Config) error {
	type EnvConfig struct {
		Address    string `env:"SERVER_ADDRESS"`
		ServerAddr string `env:"SERVER_URL"`
		FilePath   string `env:"FILE_STORAGE_PATH"`
	}

	var envConfig EnvConfig
	if err := env.Parse(&envConfig); err != nil {
		return err
	}

	// Обновляем конфигурацию только если переменные окружения заданы
	if envConfig.Address != "" {
		config.Address = envConfig.Address
	}
	if envConfig.FilePath != "" {
		config.FilePath = envConfig.FilePath
	}
	if envConfig.ServerAddr != "" {
		config.Config.ServerAddr = envConfig.ServerAddr
	}
	return nil
}

func parseFlags(config *Config) {
	flag.StringVar(&config.Address, "a", config.Address, "address for server")
	flag.StringVar(&config.Config.ServerAddr, "b", config.Config.ServerAddr, "address for link")
	flag.StringVar(&config.FilePath, "f", config.FilePath, "path to file")
	flag.Parse()
}
