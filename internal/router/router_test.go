package router

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ing-soft-2-tp1/internal/repositories"
	"testing"
)

func TestCreateRouter(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockRepo := repositories.CreateDatabase(db)

	result := CreateRouter(mockRepo)
	assert.NotNil(t, result)
}

func TestSetEnviroment(t *testing.T) {
	SetEnviroment("production")
	assert.Equal(t, gin.Mode(), gin.ReleaseMode)
}

func TestSetEnviroment2(t *testing.T) {
	SetEnviroment("Debug")
	assert.Equal(t, gin.Mode(), gin.DebugMode)
}
