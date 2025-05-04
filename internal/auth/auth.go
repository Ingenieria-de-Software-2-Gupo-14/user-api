package auth

import (
	"errors"
	"ing-soft-2-tp1/internal/config"
	e "ing-soft-2-tp1/internal/errors"
	"ing-soft-2-tp1/internal/models"
	"ing-soft-2-tp1/internal/services"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type AuthController struct {
	userRepo             services.UserService
	blockService         services.BlockedUserService
	loginAttemptsService services.LoginAttemptService
}

func NewAuthController(userRepo services.UserService) *AuthController {
	return &AuthController{
		userRepo: userRepo,
	}
}

// finishAuth finish the authentication process by generating a token and setting it in the cookie
// and redirecting to the home page
func (ac *AuthController) finishAuth(ctx *gin.Context, user models.User) {
	if user.Blocked {
		e.ErrorResponse(ctx, http.StatusForbidden, "User is blocked")
		return
	}

	token, err := GenerateToken(user.Id, user.Email, user.Name, user.Admin)
	if err != nil {
		e.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", token, 3600, "/", "", false, true)
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

// BeginAuth initiates the authentication process for the specified provider
func (ac *AuthController) BeginAuth(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "direct" {
		r := c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", c.Param("provider")))
		gothic.BeginAuthHandler(c.Writer, r)
		return
	}

	email := c.Query("email")
	password := c.Query("password")
	if email == "" || password == "" {
		e.ErrorResponse(c, http.StatusBadRequest, "Email and password are required")
		return
	}

	user, err := ac.userRepo.Login(c.Request.Context(), email, password)
	if err != nil {
		e.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	ac.finishAuth(c, *user)
}

func (ac *AuthController) CompleteAuth(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), "provider", c.Param("provider"))
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request.WithContext(ctx))
	if err != nil {
		e.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
		return
	}

	existingUser, err := ac.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		// User not found, create a new one
		if errors.Is(err, e.ErrNotFound) {
			newUser := models.CreateUserRequest{
				Name:    user.Name,
				Surname: user.LastName,
				Email:   user.Email,
			}

			user, err := ac.userRepo.CreateUser(ctx, newUser, false)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}

			existingUser = user
		} else {
			e.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
	}

	ac.finishAuth(c, *existingUser)
}

func (ac *AuthController) Logout(c *gin.Context) {
	if _, err := c.Cookie("Authorization"); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "No cookie found")
		return
	}

	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func AddAuthRoutes(r *gin.Engine, userRepo services.UserService) {
	authController := NewAuthController(userRepo)
	r.GET("/auth/:provider", authController.BeginAuth)
	r.GET("/auth/:provider/callback", authController.CompleteAuth)
}

func getAuthToken(c *gin.Context) string {
	auth, _ := c.Cookie("Authorization")
	if auth == "" {
		if parts := strings.Fields(c.GetHeader("Authorization")); len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			auth = parts[1]
		}
	}
	return auth
}

func (ac *AuthController) AuthMiddlewarefunc(ctx *gin.Context) {
	tokenStr := getAuthToken(ctx)
	claims, err := ParseToken(tokenStr)
	if err != nil {
		e.ErrorResponseWithErr(ctx, http.StatusUnauthorized, err)
		ctx.Abort()
		return
	}

	uId, _ := strconv.Atoi(claims.Subject)
	blocked, err := ac.userRepo.IsUserBlocked(ctx.Request.Context(), uId)
	if err != nil {
		e.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		ctx.Abort()
		return
	}

	if blocked {
		e.ErrorResponse(ctx, http.StatusForbidden, "User is blocked")
		// Make Cookie Expire
		ctx.SetCookie("Authorization", "", -1, "/", "", false, true)
		ctx.Abort()
		return
	}

	// Refresh token if it is about to expire
	if claims.ExpiresAt < time.Now().Add(time.Minute*5).Unix() {
		newToken, err := GenerateToken(uId, claims.Email, claims.Name, claims.Admin)
		if err != nil {
			slog.Error("Error generating new token", err)
		} else {
			ctx.SetCookie("Authorization", newToken, 3600, "/", "", false, true)
		}
	}

	ctx.Set("claims", claims)
	ctx.Next()

}
