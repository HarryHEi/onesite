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

func Success(c *gin.Context, data interface{}) {
	SuccessWithCode(c, 200, data)
}

func Created(c *gin.Context, data interface{}) {
	SuccessWithCode(c, 201, data)
}

func NoContent(c *gin.Context) {
	SuccessWithCode(c, 200, nil)
}

func BadRequest(c *gin.Context, err error) {
	FailedWithErr(c, 400, err)
}

func Unauthorized(c *gin.Context, message string) {
	FailedWithErr(c, 401, errors.New(message))
}

func PermissionDenied(c *gin.Context) {
	FailedWithErr(c, 403, errors.New("permission denied"))
}

func SuccessWithCode(c *gin.Context, code int, data interface{}) {
	resp := Response{
		Success: true,
		Code:    code,
		Msg:     "success",
		Data:    data,
	}
	c.JSON(code, resp)
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
	c.Abort()
}

// 分页查询固定响应格式
type PaginationResponse struct {
	Count int64       `json:"count"`
	Data  interface{} `json:"data"`
}
