package utils

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func NewErrorResponse(status int, title, detail, instance string) ErrorResponse {
	return ErrorResponse{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}
}

func SendErrorResponse(c *gin.Context, status int, title, detail, instance string) {
	err := NewErrorResponse(status, title, detail, instance)
	c.JSON(status, err)
}
