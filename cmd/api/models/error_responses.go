package models

const StatusMissingTitle = 1
const StatusMissingDescription = 2
const StatusUserNotFound = 3
const StatusBadRequest = 4
const StatusInternalServerError = 500
const ErrorTypeBlank = "about:blank"
const ErrorTitleMissingTitle = "Missing title"
const ErrorTitleMissingDescription = "Missing description"
const ErrorTitleUserNotFound = "User not Found"
const ErrorTitleInternalServerError = "Internal Server Error"
const ErrorTitleBadRequest = "Bad request"
const ErrorDescriptionMissingTitle = "Request sent is missing a title. please provide a title"
const ErrorDescriptionMissingDescription = "Request sent is missing a description. please provide a description"
const ErrorDescriptionUserNotFound = "User Requested doesn't exits please try another"
const ErrorDescriptionInternalServerError = "Internal server Error please try again later"
const ErrorDescriptionBadRequest = "Something's one wrong with the request please check again"

type ErrorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}
