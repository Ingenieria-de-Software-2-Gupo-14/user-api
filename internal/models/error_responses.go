package models

const (
	StatusMissingTitle                  = 1
	StatusMissingDescription            = 2
	StatusUserNotFound                  = 3
	ErrorTypeBlank                      = "about:blank"
	ErrorTitleMissingTitle              = "Missing title"
	ErrorTitleMissingDescription        = "Missing description"
	ErrorTitleUserNotFound              = "User not Found"
	ErrorTitleInternalServerError       = "Internal Server Error"
	ErrorTitleBadRequest                = "Bad request"
	ErrorDescriptionMissingTitle        = "Request sent is missing a title. please provide a title"
	ErrorDescriptionMissingDescription  = "Request sent is missing a description. please provide a description"
	ErrorDescriptionUserNotFound        = "User Requested doesn't exits please try another"
	ErrorDescriptionInternalServerError = "Internal server Error please try again later"
	ErrorDescriptionBadRequest          = "Something's gone wrong with the request please check again"
)

type ErrorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}
