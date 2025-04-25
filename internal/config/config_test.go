package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	config := LoadConfig()
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "development", config.Environment)
}
