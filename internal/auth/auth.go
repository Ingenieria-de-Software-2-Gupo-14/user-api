package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/config"
	e "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/errors"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"

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
	loginAttemptsService services.LoginAttemptService
}

func NewAuthController(userRepo services.UserService, loginAttemptsService services.LoginAttemptService) *AuthController {
	return &AuthController{
		userRepo:             userRepo,
		loginAttemptsService: loginAttemptsService,
	}
}

type authResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
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

	ac.loginAttemptsService.AddLoginAttempt(ctx, user.Id, ctx.Request.RemoteAddr, ctx.Request.UserAgent(), true)

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", token, 3600, "/", "", false, true)
	ctx.JSON(http.StatusOK, authResponse{
		User:  user,
		Token: token,
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		e.ErrorResponseWithErr(c, http.StatusBadRequest, err)
		return
	}

	user, err := ac.userRepo.Login(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		if !errors.Is(err, e.ErrNotFound) {
			ac.loginAttemptsService.AddLoginAttempt(c, user.Id, c.Request.RemoteAddr, c.Request.UserAgent(), false)
		}
		e.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	ac.finishAuth(c, *user)
}

// BeginAuth godoc
//
// @Summary      Begin authentication
// @Description  Begin authentication with the specified provider
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        provider  path      string  true  "Provider name"
// @Success      200  {object}   authResponse
// @Failure      400  {object}   e.HTTPError
// @Failure      401  {object}   e.HTTPError
// @Failure      500  {object}   e.HTTPError
// @Router       /auth/{provider} [get]
func (ac *AuthController) BeginAuth(c *gin.Context) {
	provider := c.Param("provider")
	ctx := context.WithValue(c.Request.Context(), "provider", provider)
	r := c.Request.WithContext(ctx)

	gothic.BeginAuthHandler(c.Writer, r)

}

// CompleteAuth godoc
//
// @Summary      Complete authentication
// @Description  Complete authentication with the specified provider
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        provider  path      string  true  "Provider name"
// @Success      200  {object}   authResponse
// @Failure      400  {object}   e.HTTPError
// @Failure      401  {object}   e.HTTPError
// @Failure      500  {object}   e.HTTPError
// @Router       /auth/{provider}/callback [get]
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
				e.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
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

// Logout godoc
//
// @Summary      Logout
// @Description  Logout the user by clearing the cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      307
func (ac *AuthController) Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
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
