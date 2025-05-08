package services

import (
	"net/http"
	"testing"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"

	"github.com/stretchr/testify/assert"
)

const TEST_INSTANCE = "test"

func TestCreateErrorResponse_StatusUserNotFound(t *testing.T) {
	response := CreateErrorResponse(models.StatusUserNotFound, TEST_INSTANCE)

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleUserNotFound,
		Status:   http.StatusNotFound,
		Detail:   models.ErrorDescriptionUserNotFound,
		Instance: TEST_INSTANCE,
	}

	assert.Equal(t, expectedResponse, response)
}

func TestCreateErrorResponse_StatusInternalServerError(t *testing.T) {
	response := CreateErrorResponse(http.StatusInternalServerError, TEST_INSTANCE)

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleInternalServerError,
		Status:   http.StatusInternalServerError,
		Detail:   models.ErrorDescriptionInternalServerError,
		Instance: TEST_INSTANCE,
	}

	assert.Equal(t, expectedResponse, response)
}

func TestCreateErrorResponse_StatusBadRequest(t *testing.T) {
	response := CreateErrorResponse(http.StatusBadRequest, TEST_INSTANCE)

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleBadRequest,
		Status:   http.StatusBadRequest,
		Detail:   models.ErrorDescriptionBadRequest,
		Instance: TEST_INSTANCE,
	}

	assert.Equal(t, expectedResponse, response)
}

func TestCreateErrorResponse_StatusConflict(t *testing.T) {
	response := CreateErrorResponse(http.StatusConflict, TEST_INSTANCE)

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleConflict,
		Status:   http.StatusConflict,
		Detail:   models.ErrorDescriptionConflict,
		Instance: TEST_INSTANCE,
	}

	assert.Equal(t, expectedResponse, response)
}

func TestCreateErrorResponse_StatusUnauthorized(t *testing.T) {
	response := CreateErrorResponse(http.StatusUnauthorized, TEST_INSTANCE)

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    "error",
		Status:   http.StatusUnauthorized,
		Detail:   "error",
		Instance: TEST_INSTANCE,
	}

	assert.Equal(t, expectedResponse, response)
}

func TestCreateErrorResponse_StatusForbidden(t *testing.T) {
	response := CreateErrorResponse(http.StatusForbidden, TEST_INSTANCE)

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    "error",
		Status:   http.StatusForbidden,
		Detail:   "error",
		Instance: TEST_INSTANCE,
	}

	assert.Equal(t, expectedResponse, response)
}

func TestCreateErrorResponse_UnkownStatusCode(t *testing.T) {
	response := CreateErrorResponse(0, TEST_INSTANCE)

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleInternalServerError,
		Status:   http.StatusInternalServerError,
		Detail:   models.ErrorDescriptionInternalServerError,
		Instance: TEST_INSTANCE,
	}

	assert.Equal(t, expectedResponse, response)
}
