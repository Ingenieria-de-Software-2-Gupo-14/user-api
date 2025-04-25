package auth

import (
	"ing-soft-2-tp1/internal/config"
	"ing-soft-2-tp1/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/net/context"
)

func NewAuth(config config.Config) {
	goth.UseProviders(
		google.New(config.GoogleKey, config.GoogleSecret, "http://localhost:8080/auth/google/callback"),
	)
}

func AddAuthRoutes(r *gin.Engine, userRepo *repositories.Database) {
	r.GET("/auth/:provider", func(c *gin.Context) {
		r := c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", c.Param("provider")))
		gothic.BeginAuthHandler(c.Writer, r)
	})

	r.GET("/auth/:provider/callback", func(c *gin.Context) {
		r := c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", c.Param("provider")))
		user, err := gothic.CompleteUserAuth(c.Writer, r)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		existingUser, err := userRepo.GetUserByEmail(c.Request.Context(), user.Email)
		if err != nil {
			// User not found, create a new one
		}

		if existingUser.BlockedUser {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is blocked"})
			return
		}

		token, err := GenerateToken(*existingUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.SetCookie("Authorization", token, 3600, "/", "", false, true)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})
}

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

		jwt, err := ParseToken(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("jwt_info", jwt)
		ctx.Next()
	}
}
