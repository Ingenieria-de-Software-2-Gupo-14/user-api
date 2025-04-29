package auth

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ing-soft-2-tp1/internal/config"
	"ing-soft-2-tp1/internal/repositories"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddAuthRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, _ := sqlmock.New()

	r := gin.Default()
	repo := repositories.CreateUserRepo(db)

	AddAuthRoutes(r, repo)
	assert.True(t, routeExists(r, "GET", "/auth/:provider"))
	assert.True(t, routeExists(r, "GET", "/auth/:provider/callback"))
}

func TestAddAuthRoutes_ProviderGet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, _ := sqlmock.New()

	r := gin.Default()
	repo := repositories.CreateUserRepo(db)

	AddAuthRoutes(r, repo)
	req, _ := http.NewRequest("GET", "/auth/google", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestAddAuthRoutes_Callback_InternalServerError_WrongProvider(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, _ := sqlmock.New()

	r := gin.Default()
	repo := repositories.CreateUserRepo(db)

	AddAuthRoutes(r, repo)
	req, _ := http.NewRequest("GET", "/auth/google/callback", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestNewAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testcCongif := config.Config{
		Host:         "",
		Port:         "",
		Environment:  "",
		DatabaseURL:  "",
		GoogleKey:    "",
		GoogleSecret: "",
		JWTSecret:    "",
	}

	NewAuth(testcCongif)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, _ := sqlmock.New()

	repo := repositories.CreateUserRepo(db)

	testFunc := AuthMiddleware(repo)
	assert.NotNil(t, testFunc)
}

func routeExists(engine *gin.Engine, method, path string) bool {
	for _, r := range engine.Routes() {
		if r.Method == method && r.Path == path {
			return true
		}
	}
	return false
}
