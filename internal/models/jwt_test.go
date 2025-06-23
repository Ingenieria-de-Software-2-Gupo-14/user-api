package models

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
)

func TestGetJWTSecret(t *testing.T) {
	original := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", original)

	t.Run("returns env value when set", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "mysecret")
		secret := GetJWTSecret()
		require.Equal(t, "mysecret", secret)
	})

	t.Run("returns default value when env not set", func(t *testing.T) {
		os.Unsetenv("JWT_SECRET")
		secret := GetJWTSecret()
		require.Equal(t, "secret", secret)
	})
}

func TestParseInvalidToken(t *testing.T) {
	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"

	token, err := GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	claims, err := ParseToken(token + "Invalidsignature")
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.ErrorIs(t, err, ErrJWTValidation)
}

func TestGenerateToken(t *testing.T) {
	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"

	token, err := GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	parseToken, err := ParseToken(token)
	assert.NoError(t, err)

	assert.Equal(t, parseToken.Role, role)
	assert.Equal(t, parseToken.Email, email)
	assert.Equal(t, parseToken.Name, name)
	tokenUserId, err := strconv.Atoi(parseToken.Subject)
	assert.NoError(t, err)
	assert.Equal(t, tokenUserId, userId)
}
