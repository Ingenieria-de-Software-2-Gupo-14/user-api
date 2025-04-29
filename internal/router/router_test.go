package router

import (
	"github.com/DATA-DOG/go-sqlmock"
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

func TestCreateRouter(t *testing.T) {
	db, _, err := sqlmock.New()

	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)

	router := CreateRouter(db)

	assert.NotNil(t, router)
}
