package router

import (
	"testing"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetEnviroment(t *testing.T) {
	SetEnviroment("production")
	assert.Equal(t, gin.Mode(), gin.ReleaseMode)
}

func TestSetEnviroment2(t *testing.T) {
	SetEnviroment("Debug")
	assert.Equal(t, gin.Mode(), gin.DebugMode)
}

func TestCreateRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router, err := CreateRouter(config.Config{})
	assert.Error(t, err)
	assert.Nil(t, router)
}
