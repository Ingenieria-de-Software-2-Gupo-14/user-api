package middleware

import (
	"ing-soft-2-tp1/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getAuthToken(c *gin.Context) string {
	auth, _ := c.Cookie("Authorization")
	if auth == "" {
		auth = c.Request.Header.Get("Authorization")
	}
	return auth
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := getAuthToken(ctx)

		jwt, err := models.ParseToken(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("jwt_info", jwt)
		ctx.Next()
	}
}
