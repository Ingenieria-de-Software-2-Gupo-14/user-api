package services

import (
	"net/http"

	. "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
)

// CreateErrorResponse creates and returns an error response struct according to the status code, in case of not recognizing the status code and internal server error response is created
func CreateErrorResponse(statusCode int, instance string) ErrorResponse {
	switch statusCode {
	case StatusUserNotFound:
		return ErrorResponse{Type: ErrorTypeBlank, Title: ErrorTitleUserNotFound, Status: http.StatusNotFound, Detail: ErrorDescriptionUserNotFound, Instance: instance}
	case http.StatusInternalServerError:
		return ErrorResponse{Type: ErrorTypeBlank, Title: ErrorTitleInternalServerError, Status: http.StatusInternalServerError, Detail: ErrorDescriptionInternalServerError, Instance: instance}
	case http.StatusBadRequest:
		return ErrorResponse{Type: ErrorTypeBlank, Title: ErrorTitleBadRequest, Status: http.StatusBadRequest, Detail: ErrorDescriptionBadRequest, Instance: instance}
	case http.StatusConflict:
		return ErrorResponse{Type: ErrorTypeBlank, Title: ErrorTitleConflict, Status: http.StatusConflict, Detail: ErrorDescriptionConflict, Instance: instance}
	case http.StatusUnauthorized:
		return ErrorResponse{Type: ErrorTypeBlank, Title: "error", Status: http.StatusUnauthorized, Detail: "error", Instance: instance}
	case http.StatusForbidden:
		return ErrorResponse{Type: ErrorTypeBlank, Title: "error", Status: http.StatusForbidden, Detail: "error", Instance: instance}
	}
	return ErrorResponse{Type: ErrorTypeBlank, Title: ErrorTitleInternalServerError, Status: http.StatusInternalServerError, Detail: ErrorDescriptionInternalServerError, Instance: instance}
}
