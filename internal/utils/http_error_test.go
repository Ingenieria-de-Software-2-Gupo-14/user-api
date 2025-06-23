package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorResponseWithCustomError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	err := errors.New("something went wrong")
	ErrorResponseWithErr(c, http.StatusBadRequest, err)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{
		"code": 400,
		"title": "Bad Request",
		"error": "something went wrong"
	}`, w.Body.String())
}

func TestErrorResponseWithNilError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ErrorResponseWithErr(c, http.StatusInternalServerError, nil)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.Contains(t, w.Body.String(), "unknown error") // assumes ErrUnkwon.Error() == "unknown error"
}

func TestErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ErrorResponse(c, http.StatusForbidden, "you are not allowed")

	require.Equal(t, http.StatusForbidden, w.Code)
	require.JSONEq(t, `{
		"code": 403,
		"title": "Forbidden",
		"error": "you are not allowed"
	}`, w.Body.String())
}
