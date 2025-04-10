// Description: Config for running the server
package config

import (
	"encoding/json"
	"flag"
	"os"
	"shorter/internal/handlers"

	"github.com/caarlos0/env/v11"
)

// Config struct
type Config struct {
	Address     string `env:"SERVER_ADDRESS" json:"base_url"`             // Address for server
	FilePath    string `env:"FILE_STORAGE_PATH" json:"file_storage_path"` // Optional if you want to save in file
	ServerAddr  string `env:"SERVER_URL" json:"server_address"`           // Host for returned link with short url
	DatabaseURL string `env:"DATABASE_DSN" json:"database_dsn"`           // Optional if you want to save in database
	HTTPS       bool   `env:"ENABLE_HTTPS" json:"enable_https"`           // Optional if you want to use https
	CfgPath     string `env:"CONFIG"`                                     // Optional if you want to config.json
	Config      *handlers.Config
}

// NewConfig creates a new config
func NewConfig() (*Config, error) {
	config := &Config{
		Address:  ":8080",
		FilePath: "",
		HTTPS:    false,
		Config: &handlers.Config{
			ServerAddr: "http://localhost:8080",
		},
	}

	parseFlags(config)
	if err := parseEnv(config); err != nil {
		return nil, err
	}
	if err := parseJSON(config); err != nil {
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
	if envConfig.HTTPS {
		config.HTTPS = true
	}
	if envConfig.CfgPath != "" {
		config.CfgPath = envConfig.CfgPath
	}
	return nil
}

func parseFlags(config *Config) {
	flag.StringVar(&config.Address, "a", config.Address, "address for server")
	flag.StringVar(&config.Config.ServerAddr, "b", config.Config.ServerAddr, "address for link")
	flag.StringVar(&config.FilePath, "f", config.FilePath, "path to file")
	flag.StringVar(&config.Config.DatabaseURL, "d", config.Config.DatabaseURL, "path to file")
	flag.BoolVar(&config.HTTPS, "s", config.HTTPS, "enable HTTPS (default: false)")
	flag.StringVar(&config.CfgPath, "c", config.CfgPath, "Config path default empty")
	flag.StringVar(&config.CfgPath, "config", config.CfgPath, "Config path default empty long version")
	flag.Parse()
}

func parseJSON(config *Config) error {
	if config.CfgPath == "" {
		return nil
	}
	file, err := os.Open(config.CfgPath)
	if err != nil {
		return err
	}
	defer file.Close()
	var jsonConf Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonConf); err != nil {
		return err
	}
	if jsonConf.ServerAddr != "" && config.Config.ServerAddr == "" {
		config.Config.ServerAddr = jsonConf.ServerAddr
	}
	if jsonConf.Address != "" && config.Address == "" {
		config.Address = jsonConf.Address
	}
	if jsonConf.FilePath != "" && config.FilePath == "" {
		config.FilePath = jsonConf.FilePath
	}
	if jsonConf.DatabaseURL != "" && config.DatabaseURL == "" {
		config.DatabaseURL = jsonConf.DatabaseURL
	}
	if jsonConf.HTTPS && !config.HTTPS {
		config.HTTPS = jsonConf.HTTPS
	}
	return nil
}
