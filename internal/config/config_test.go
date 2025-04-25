package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := LoadConfig()
	assert.Equal(t, os.Getenv("PORT"), config.Port)
	assert.Equal(t, os.Getenv("HOST"), config.Host)
	assert.Equal(t, os.Getenv("ENVIRONMENT"), config.Environment)
}
