package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Code  int    `json:"code"`
	Title string `json:"title"`
	Error string `json:"error"`
}

var ErrUnkwon = errors.New("unknown error")

func ErrorResponseWithErr(ctx *gin.Context, code int, err error) {
	if err == nil {
		err = ErrUnkwon
	}

	ctx.JSON(code, HTTPError{
		Code:  code,
		Title: http.StatusText(code),
		Error: err.Error(),
	})
}

func ErrorResponse(ctx *gin.Context, code int, errorMessage string) {
	ctx.JSON(code, HTTPError{
		Code:  code,
		Title: http.StatusText(code),
		Error: errorMessage,
	})
}
