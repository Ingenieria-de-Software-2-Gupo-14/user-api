package auth

import (
	"ing-soft-2-tp1/internal/models"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {

	tokenStr, err := GenerateToken(models.User{})
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	info, err := ParseToken(tokenStr)
	assert.NoError(t, err)

	assert.Equal(t, 1, (info.UserId))
	assert.Equal(t, false, info.Admin)
	assert.True(t, info.Exp > info.Iat)
	assert.True(t, info.Exp > time.Now().Unix())
}

func TestParseToken_Fail(t *testing.T) {

	result, err := ParseToken("wawa")

	assert.NotNil(t, err)
	assert.Equal(t, JWTInfo{}, result)

}

func TestNewJWTInfoFromClaims_Error(t *testing.T) {
	claims := jwt.MapClaims{
		"errorTest": "error",
	}

	_, err := NewJWTInfoFromClaims(claims)
	assert.NotNil(t, err)
}
