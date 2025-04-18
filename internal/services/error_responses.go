package services

import (
	. "ing-soft-2-tp1/internal/models"
	"net/http"
)

// CreateErrorResponse creates and returns an error response struct according to the status code, in case of not recognizing the status code and internal server error response is created
func CreateErrorResponse(statusCode int, instance string) ErrorResponse {
	switch statusCode {
	case StatusMissingTitle:
		return ErrorResponse{Type: ErrorTypeBlank, Title: ErrorTitleMissingTitle, Status: http.StatusBadRequest, Detail: ErrorDescriptionMissingTitle, Instance: instance}
	case StatusMissingDescription:
		return ErrorResponse{Type: ErrorTypeBlank, Title: ErrorTitleMissingDescription, Status: http.StatusBadRequest, Detail: ErrorDescriptionMissingDescription, Instance: instance}
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
	}
	return ErrorResponse{Type: ErrorTypeBlank, Title: ErrorTitleInternalServerError, Status: http.StatusInternalServerError, Detail: ErrorDescriptionInternalServerError, Instance: instance}
}
