package middleware

import (
	"database/sql"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func performRequestWithToken(r http.Handler, token string, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := 123

	mockService := services.NewMockUserService(t)
	mockService.EXPECT().IsUserBlocked(mock.Anything, userID).Return(false, nil)

	r := gin.New()
	r.Use(AuthMiddleware(mockService))
	r.GET("/", func(c *gin.Context) {
		claims, exists := c.Get("claims")
		assert.True(t, exists)
		assert.NotNil(t, claims)
		c.Status(http.StatusOK)
	})

	email := "test@test.com"
	name := "test"
	role := "user"
	token, err := models.GenerateToken(userID, email, name, role)
	assert.NoError(t, err)
	w := performRequestWithToken(r, token, "/")

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuthMiddleware_BlockedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := 1
	email := "test@test.com"
	name := "test"
	role := "user"

	token, err := models.GenerateToken(userID, email, name, role)
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)
	mockService.EXPECT().IsUserBlocked(mock.Anything, userID).Return(true, nil)

	r := gin.New()
	r.Use(AuthMiddleware(mockService))
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, token, "/")

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	mockService := services.NewMockUserService(t)

	r.Use(AuthMiddleware(mockService))
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, "invalid token", "/")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_IsUserBlockedError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := 10
	email := "test@test.com"
	name := "test"
	role := "user"

	token, err := models.GenerateToken(userID, email, name, role)
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)
	mockService.EXPECT().IsUserBlocked(mock.Anything, userID).Return(false, sql.ErrConnDone)

	r := gin.New()
	r.Use(AuthMiddleware(mockService))
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, token, "/")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuthMiddleware_TokenRefresh(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := 1
	email := "test@example.com"
	name := "test"
	role := "user"

	claims := models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.Itoa(userID),
			Issuer:    "user-api",
			ExpiresAt: time.Now().Add(2 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Email: email,
		Name:  name,
		Role:  role,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(models.GetJWTSecret()))
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)
	mockService.EXPECT().IsUserBlocked(mock.Anything, userID).Return(false, nil)

	r := gin.New()
	r.Use(AuthMiddleware(mockService))
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, token, "/")

	assert.Equal(t, http.StatusOK, w.Code)

	found := false
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "Authorization" {
			found = true
			assert.Greater(t, cookie.MaxAge, 0)
		}
	}
	assert.True(t, found, "Expected refreshed token cookie")
}

func TestAdminOnlyMiddleware_AdminAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := 1
	email := "admin@test.com"
	name := "admin"
	role := "admin"

	token, err := models.GenerateToken(userID, email, name, role)
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)

	r := gin.New()
	r.Use(AdminOnlyMiddleware(mockService))
	r.GET("/", func(c *gin.Context) {
		claims, exists := c.Get("claims")
		assert.True(t, exists)
		assert.NotNil(t, claims)
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, token, "/")

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminOnlyMiddleware_NonAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := 2
	email := "user@test.com"
	name := "user"
	role := "user"

	token, err := models.GenerateToken(userID, email, name, role)
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)

	r := gin.New()
	r.Use(AdminOnlyMiddleware(mockService))
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, token, "/")

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAdminOnlyMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := services.NewMockUserService(t)

	r := gin.New()
	r.Use(AdminOnlyMiddleware(mockService))
	r.GET("/", func(c *gin.Context) {
		claims, exists := c.Get("claims")
		assert.True(t, exists)
		assert.NotNil(t, claims)
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, "invalid token", "/")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminOnlyMiddleware_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := 1
	email := "test@example.com"
	name := "test"
	role := "admin"

	claims := models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.Itoa(userID),
			Issuer:    "user-api",
			ExpiresAt: time.Now().Add(2 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Email: email,
		Name:  name,
		Role:  role,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(models.GetJWTSecret()))
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)

	r := gin.New()
	r.Use(AdminOnlyMiddleware(mockService))
	r.GET("/", func(c *gin.Context) {
		claims, exists := c.Get("claims")
		assert.True(t, exists)
		assert.NotNil(t, claims)
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, token, "/")

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserOrAdminMiddleware_NormalUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := 2
	email := "user@test.com"
	name := "user"
	role := "user"

	token, err := models.GenerateToken(userID, email, name, role)
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)
	mockService.EXPECT().IsUserBlocked(mock.Anything, userID).Return(false, nil)

	r := gin.New()
	r.Use(UserOrAdminMiddleware(mockService))
	r.GET("/users/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, token, "/users/2")

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserOrAdminMiddleware_BlockedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := 1
	token, err := models.GenerateToken(userID, "user@test.com", "User", "user")
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)
	mockService.EXPECT().IsUserBlocked(mock.Anything, userID).Return(true, nil)

	r := gin.New()
	r.Use(UserOrAdminMiddleware(mockService))
	r.GET("/users/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, token, "/users/1")

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUserOrAdminMiddleware_ForbiddenAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := 1
	token, err := models.GenerateToken(userID, "user@test.com", "User", "user")
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)
	mockService.EXPECT().IsUserBlocked(mock.Anything, userID).Return(false, nil)

	r := gin.New()
	r.Use(UserOrAdminMiddleware(mockService))
	r.GET("/users/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, token, "/users/99")

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUserOrAdminMiddleware_AdminAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	adminID := 1
	token, err := models.GenerateToken(adminID, "admin@test.com", "Admin", "admin")
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)
	mockService.EXPECT().IsUserBlocked(mock.Anything, adminID).Return(false, nil)

	r := gin.New()
	r.Use(UserOrAdminMiddleware(mockService))
	r.GET("/users/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, token, "/users/999")

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserOrAdminMiddleware_MissingPathParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := 1
	token, err := models.GenerateToken(userID, "user@test.com", "User", "user")
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)
	mockService.EXPECT().IsUserBlocked(mock.Anything, userID).Return(false, nil)

	r := gin.New()
	r.Use(UserOrAdminMiddleware(mockService))
	r.GET("/users", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// This route does not include /:id
	w := performRequestWithToken(r, token, "/users")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserOrAdminMiddleware_InvalidPathParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := 1
	token, err := models.GenerateToken(userID, "user@test.com", "User", "user")
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)
	mockService.EXPECT().IsUserBlocked(mock.Anything, userID).Return(false, nil)

	r := gin.New()
	r.Use(UserOrAdminMiddleware(mockService))
	r.GET("/users/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, token, "/users/abc")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserOrAdminMiddleware_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := 1
	email := "test@example.com"
	name := "test"
	role := "user"

	claims := models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.Itoa(userID),
			Issuer:    "user-api",
			ExpiresAt: time.Now().Add(2 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Email: email,
		Name:  name,
		Role:  role,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(models.GetJWTSecret()))
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t)
	mockService.EXPECT().IsUserBlocked(mock.Anything, userID).Return(false, nil)

	r := gin.New()
	r.Use(UserOrAdminMiddleware(mockService))
	r.GET("/users/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := performRequestWithToken(r, token, "/users/1")

	assert.Equal(t, http.StatusOK, w.Code)

	found := false
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "Authorization" {
			found = true
			assert.Greater(t, cookie.MaxAge, 0)
		}
	}
	assert.True(t, found, "Expected refreshed token cookie")
}
