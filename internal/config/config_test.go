package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	config := LoadConfig()
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "development", config.Environment)
}

func TestCreateDatabase_WrongUrl(t *testing.T) {
	testConfig := Config{
		Host:         "",
		Port:         "",
		Environment:  "",
		DatabaseURL:  "",
		GoogleKey:    "",
		GoogleSecret: "",
		JWTSecret:    "",
	}

	_, err := testConfig.CreateDatabase()
	assert.NotNil(t, err)
	println(err)
}

func TestLoadConfig2(t *testing.T) {
	err := os.Setenv("DATABASE_URL", "a")
	if err != nil {
		return
	}
	config := LoadConfig()
	assert.Equal(t, config.DatabaseURL, "a")
}

func TestLoadConfig3(t *testing.T) {
	err := os.Setenv("HOST", "a")
	if err != nil {
		return
	}
	config := LoadConfig()
	assert.Equal(t, config.Host, "a")
}
