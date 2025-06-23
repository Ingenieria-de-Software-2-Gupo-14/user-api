package middleware

import (
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestAdminOnlyMiddleware_AdminAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := 1
	email := "admin@test.com"
	name := "admin"
	role := "admin"

	token, err := models.GenerateToken(userID, email, name, role)
	assert.NoError(t, err)

	mockService := services.NewMockUserService(t) // not used, but required for signature

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
