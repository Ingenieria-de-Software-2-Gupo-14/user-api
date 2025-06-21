package router

import (
	"os"
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
	os.Setenv("TESTING", "true")
	gin.SetMode(gin.TestMode)

	router, err := CreateRouter(config.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, router)
	os.Setenv("TESTING", "")
}
