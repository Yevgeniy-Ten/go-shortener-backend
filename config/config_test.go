package config

import (
	"os"
	"shorter/internal/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig_DefaultValues(t *testing.T) {
	os.Clearenv()

	cfg, err := NewConfig()
	require.NoError(t, err)

	assert.Equal(t, ":8080", cfg.Address)
	assert.Equal(t, "", cfg.FilePath)
	assert.Equal(t, "http://localhost:8080", cfg.Config.ServerAddr)
	assert.Equal(t, "", cfg.Config.DatabaseURL)
}
func TestParseEnv(t *testing.T) {
	os.Setenv("SERVER_ADDRESS", ":9090")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/file")
	os.Setenv("SERVER_URL", "http://localhost:9091")
	os.Setenv("DATABASE_DSN", "postgres://localhost:5432/testdb")

	config := &Config{
		Config: &handlers.Config{},
	}

	err := parseEnv(config)
	require.NoError(t, err)

	assert.Equal(t, ":9090", config.Address)
	assert.Equal(t, "/tmp/file", config.FilePath)
	assert.Equal(t, "http://localhost:9091", config.Config.ServerAddr)
	assert.Equal(t, "postgres://localhost:5432/testdb", config.Config.DatabaseURL)
}
