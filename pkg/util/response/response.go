package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 基础响应结构体
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

const (
	SUCCESS = 0
	ERROR   = 500
)

// Result 返回统一规范的 JSON 数据
func Result(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Data:    data,
		Message: msg,
	})
}

// Success 成功返回
func Success(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, data, "success")
}

// SuccessWithMessage 自定义成功信息返回
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	Result(c, SUCCESS, data, message)
}

// Fail 失败返回
func Fail(c *gin.Context, code int, msg string) {
	Result(c, code, nil, msg)
}

// FailWithMessage 简单错误返回 (默认 500)
func FailWithMessage(c *gin.Context, msg string) {
	Result(c, ERROR, nil, msg)
}
