package auth

import (
	"ing-soft-2-tp1/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Authenticated")
	})

	token, err := GenerateToken(models.User{})
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
