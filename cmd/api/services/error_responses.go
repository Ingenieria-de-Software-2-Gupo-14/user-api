package services

import (
	. "ing-soft-2-tp1/cmd/api/models"
	"net/http"
)

// CreateErrorResponse creates and returns an error response struct according to the status code, in case of not recognizing the status code and internal server error response is created
func CreateErrorResponse(statusCode int, instance string) ErrorResponse {
	switch statusCode {
	case StatusMissingTitle:
		return ErrorResponse{ErrorTypeBlank, ErrorTitleMissingTitle, http.StatusBadRequest, ErrorDescriptionMissingTitle, instance}
	case StatusMissingDescription:
		return ErrorResponse{ErrorTypeBlank, ErrorTitleMissingDescription, http.StatusBadRequest, ErrorDescriptionMissingDescription, instance}
	case StatusUserNotFound:
		return ErrorResponse{ErrorTypeBlank, ErrorTitleUserNotFound, http.StatusNotFound, ErrorDescriptionUserNotFound, instance}
	case StatusInternalServerError:
		return ErrorResponse{ErrorTypeBlank, ErrorTitleInternalServerError, http.StatusInternalServerError, ErrorDescriptionInternalServerError, instance}
	case StatusBadRequest:
		return ErrorResponse{ErrorTypeBlank, ErrorTitleBadRequest, http.StatusBadRequest, ErrorDescriptionBadRequest, instance}
	}
	return ErrorResponse{ErrorTypeBlank, ErrorTitleInternalServerError, http.StatusInternalServerError, ErrorDescriptionInternalServerError, instance}
}
