package auth

import (
	"errors"
	"ing-soft-2-tp1/internal/config"
	e "ing-soft-2-tp1/internal/errors"
	"ing-soft-2-tp1/internal/models"
	"ing-soft-2-tp1/internal/services"
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

func AddAuthRoutes(r *gin.Engine, userRepo services.UserRepository) {
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
			if errors.Is(err, e.ErrNotFound) {
				newUser := &models.User{
					Username: user.NickName,
					Name:     user.Name,
					Email:    user.Email,
				}

				id, err := userRepo.AddUser(c.Request.Context(), newUser)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
					return
				}

				newUser.Id = id
				existingUser = newUser
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
				return
			}
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

func AuthMiddleware(userRepo services.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := getAuthToken(ctx)

		jwt, err := ParseToken(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": err.Error()})
			ctx.Abort()
			return
		}

		user, err := userRepo.GetUserByEmail(ctx.Request.Context(), jwt.Email)
		if err != nil {
			if errors.Is(err, e.ErrNotFound) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "User not found"})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
			ctx.Abort()
			return
		}

		if user.BlockedUser {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "User is blocked"})
			ctx.Abort()
			return
		}

		ctx.Set("jwt_info", jwt)
		ctx.Next()
	}
}
