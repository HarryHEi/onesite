package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

func NoContent(c *gin.Context) {
	Success(c, nil)
}

func BadRequest(c *gin.Context, err error) {
	FailedWithErr(c, 400, err)
}

func Unauthorized(c *gin.Context, message string) {
	FailedWithErr(c, 401, errors.New(message))
}

func Success(c *gin.Context, data interface{}) {
	resp := Response{
		Success: true,
		Code:    200,
		Msg:     "success",
		Data:    data,
	}
	c.JSON(200, resp)
	return
}

func FailedWithErr(c *gin.Context, code int, err error) {
	resp := Response{
		Success: false,
		Code:    code,
		Msg:     err.Error(),
		Data:    nil,
	}
	c.JSON(code, resp)
}
