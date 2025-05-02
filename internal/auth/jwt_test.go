package auth

import (
	"ing-soft-2-tp1/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	tokenStr, err := GenerateToken(models.User{})
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	ou, err := ParseToken(tokenStr)

	assert.Equal(t, ou, Claims{})
	assert.NoError(t, err)
}
