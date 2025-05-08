package auth

import (
	"testing"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"

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
