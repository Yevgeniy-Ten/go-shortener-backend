// Description: Config for running the server
package config

import (
	"flag"
	"shorter/internal/handlers"

	"github.com/caarlos0/env/v11"
)

// Config struct
type Config struct {
	Address     string `env:"SERVER_ADDRESS"`    // Address for server
	FilePath    string `env:"FILE_STORAGE_PATH"` // Optional if you want to save in file
	ServerAddr  string `env:"SERVER_URL"`        // Host for returned link with short url
	DatabaseURL string `env:"DATABASE_DSN"`      // Optional if you want to save in database
	HTTPs       bool   `env:"ENABLE_HTTPS"`      // Optional if you want to use https
	Config      *handlers.Config
}

// NewConfig creates a new config
func NewConfig() (*Config, error) {
	config := &Config{
		Address:  ":8080",
		FilePath: "",
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
	var envConfig Config
	if err := env.Parse(&envConfig); err != nil {
		return err
	}

	if envConfig.Address != "" {
		config.Address = envConfig.Address
	}
	if envConfig.FilePath != "" {
		config.FilePath = envConfig.FilePath
	}
	if envConfig.ServerAddr != "" {
		config.Config.ServerAddr = envConfig.ServerAddr
	}
	if envConfig.DatabaseURL != "" {
		config.Config.DatabaseURL = envConfig.DatabaseURL
	}
	if envConfig.HTTPs {
		config.HTTPs = true
	}
	return nil
}

func parseFlags(config *Config) {
	flag.StringVar(&config.Address, "a", config.Address, "address for server")
	flag.StringVar(&config.Config.ServerAddr, "b", config.Config.ServerAddr, "address for link")
	flag.StringVar(&config.FilePath, "f", config.FilePath, "path to file")
	flag.StringVar(&config.Config.DatabaseURL, "d", config.Config.DatabaseURL, "path to file")
	flag.BoolVar(&config.HTTPs, "s", config.HTTPs, "enable HTTPS (default: false)")
	flag.Parse()
}
