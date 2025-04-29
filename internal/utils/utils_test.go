package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	result, err := HashPassword("test")
	if err != nil {
		assert.Fail(t, "Something went wrong while Hashing the password")
	}
	assert.True(t, result != "")
}

func TestCompareHashPassword(t *testing.T) {
	result, errHash := HashPassword("test")
	if errHash != nil {
		assert.Fail(t, "Something went wrong while Hashing the password")
	}
	errCompare := CompareHashPassword(result, "test")
	assert.NoError(t, errCompare)
}

func TestCompareHashPassword2(t *testing.T) {
	result, errHash := HashPassword("test")
	if errHash != nil {
		assert.Fail(t, "Something went wrong while Hashing the password")
	}
	errCompare := CompareHashPassword(result, "example")
	assert.Error(t, errCompare)
}
