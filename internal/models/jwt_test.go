package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {

	tokenStr, err := GenerateToken(1, false)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	info, err := ParseToken(tokenStr)
	assert.NoError(t, err)

	assert.Equal(t, 1, (info.UserId))
	assert.Equal(t, false, info.Admin)
	assert.True(t, info.Exp > info.Iat)
	assert.True(t, info.Exp > time.Now().Unix())
}
