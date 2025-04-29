package router

import (
	"testing"

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
